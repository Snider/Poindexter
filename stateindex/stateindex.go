// Package stateindex is a sidecar index over State container files.
//
// A State `.kv` container holds a JSON header and a binary tail. Once one file
// can hold several State segments, or a session spans several files, something
// has to answer "where does the State behind this index URI actually live?"
// without opening containers to find out. That is this package: a small JSON
// document listing entries, each naming a file and the byte range of the
// payload inside it.
//
// Nothing here reads a container. Entries carry the offset and length recorded
// when the container was written, so resolving an index URI to a byte range
// costs a map lookup — the binary State payload is never touched, and the
// files an index refers to need not even be present.
//
//	idx := stateindex.New()
//	_ = idx.Append(stateindex.Entry{
//	    ID:            "session-1-fold-1",
//	    Path:          "session-1.kv",
//	    IndexURI:      "mlx://state-ramp/fold/1/folded/index",
//	    TokenCount:    206,
//	    PayloadOffset: 1234,
//	    PayloadBytes:  80511040,
//	})
//	entry, ok := idx.Lookup("mlx://state-ramp/fold/1/folded/index")
package stateindex

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// DefaultKind is the value New puts in an index's Kind field.
//
// Kind is a caller-owned label, not something this package enforces — Decode
// accepts any value, including an empty one. Change it freely if the consuming
// project has agreed a different tag.
const DefaultKind = "lthn/state-index"

var (
	// ErrMissingIndexURI is returned when an entry has no index URI, which is
	// the key everything here is addressed by.
	ErrMissingIndexURI = errors.New("stateindex: entry has no index_uri")
	// ErrMissingPath is returned when an entry names no container file.
	ErrMissingPath = errors.New("stateindex: entry has no path")
	// ErrNegativeRange is returned when an entry's payload range is negative.
	ErrNegativeRange = errors.New("stateindex: entry has a negative payload range")
	// ErrDuplicateIndexURI is returned when appending an index URI the index
	// already carries. Use Set to replace an entry deliberately.
	ErrDuplicateIndexURI = errors.New("stateindex: index_uri already present")
)

// Entry locates one State segment: which container file holds it, and where
// its payload starts and ends inside that file.
type Entry struct {
	// ID is the caller's name for the segment, for display and debugging.
	ID string `json:"id"`
	// Path is the container file holding this segment. Entries in one index
	// may name the same file or many different ones.
	Path string `json:"path"`
	// IndexURI is the key the segment is looked up by, and is unique within
	// an index.
	IndexURI string `json:"index_uri"`
	// TokenCount is how many tokens the segment represents, when known.
	TokenCount int64 `json:"token_count,omitempty"`
	// PayloadOffset is the byte offset of the segment's payload within Path —
	// the value an mmap or pread starts from.
	PayloadOffset int64 `json:"payload_offset"`
	// PayloadBytes is the length of the segment's payload in bytes.
	PayloadBytes int64 `json:"payload_bytes"`
}

// Index is the sidecar document: a kind tag and an ordered list of entries.
//
// The zero value is usable but carries no Kind; prefer New.
type Index struct {
	Kind   string  `json:"kind"`
	States []Entry `json:"states"`
}

// New returns an empty index tagged with DefaultKind.
func New() *Index {
	return &Index{Kind: DefaultKind}
}

// Append adds an entry, rejecting one that is unaddressable (no index URI, no
// path, a negative range) or that duplicates an index URI already present.
//
// Order is preserved: an index reads back in the order it was built.
func (i *Index) Append(e Entry) error {
	if err := e.validate(); err != nil {
		return err
	}
	if _, ok := i.Lookup(e.IndexURI); ok {
		return fmt.Errorf("%w: %s", ErrDuplicateIndexURI, e.IndexURI)
	}
	i.States = append(i.States, e)
	return nil
}

// Set adds an entry, replacing any existing entry with the same index URI in
// place. Use it when re-packing a segment; use Append when a duplicate is a
// bug worth hearing about.
func (i *Index) Set(e Entry) error {
	if err := e.validate(); err != nil {
		return err
	}
	for n := range i.States {
		if i.States[n].IndexURI == e.IndexURI {
			i.States[n] = e
			return nil
		}
	}
	i.States = append(i.States, e)
	return nil
}

// Lookup returns the entry for an index URI.
func (i *Index) Lookup(indexURI string) (Entry, bool) {
	for _, e := range i.States {
		if e.IndexURI == indexURI {
			return e, true
		}
	}
	return Entry{}, false
}

// Remove drops the entry for an index URI, reporting whether one was there.
func (i *Index) Remove(indexURI string) bool {
	for n, e := range i.States {
		if e.IndexURI == indexURI {
			i.States = append(i.States[:n], i.States[n+1:]...)
			return true
		}
	}
	return false
}

// ByPath returns every entry stored in one container file, in index order.
// An index may point at a single `.kv` file or at many.
func (i *Index) ByPath(path string) []Entry {
	var out []Entry
	for _, e := range i.States {
		if e.Path == path {
			out = append(out, e)
		}
	}
	return out
}

// Paths returns the distinct container files this index refers to, in the
// order they first appear — the set of files a consumer needs to open.
func (i *Index) Paths() []string {
	seen := make(map[string]struct{}, len(i.States))
	out := make([]string, 0, len(i.States))
	for _, e := range i.States {
		if _, ok := seen[e.Path]; ok {
			continue
		}
		seen[e.Path] = struct{}{}
		out = append(out, e.Path)
	}
	return out
}

// Len returns the number of entries.
func (i *Index) Len() int { return len(i.States) }

// Encode writes the index as JSON.
func (i *Index) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(i)
}

// Decode reads a JSON index. Entries are validated, so a sidecar naming a
// segment nothing can address fails here rather than at resolve time.
//
// Kind is not checked: it is a caller-owned label.
func Decode(r io.Reader) (*Index, error) {
	if r == nil {
		return nil, errors.New("stateindex: reader cannot be nil")
	}
	var idx Index
	if err := json.NewDecoder(r).Decode(&idx); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{}, len(idx.States))
	for n, e := range idx.States {
		if err := e.validate(); err != nil {
			return nil, fmt.Errorf("stateindex: entry %d: %w", n, err)
		}
		if _, ok := seen[e.IndexURI]; ok {
			return nil, fmt.Errorf("stateindex: entry %d: %w: %s", n, ErrDuplicateIndexURI, e.IndexURI)
		}
		seen[e.IndexURI] = struct{}{}
	}
	return &idx, nil
}

// validate rejects entries that could not be resolved to a byte range.
func (e Entry) validate() error {
	if e.IndexURI == "" {
		return ErrMissingIndexURI
	}
	if e.Path == "" {
		return ErrMissingPath
	}
	if e.PayloadOffset < 0 || e.PayloadBytes < 0 {
		return fmt.Errorf("%w: offset %d, bytes %d", ErrNegativeRange, e.PayloadOffset, e.PayloadBytes)
	}
	return nil
}
