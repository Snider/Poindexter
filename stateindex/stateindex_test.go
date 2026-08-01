package stateindex

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleEntry() Entry {
	return Entry{
		ID:            "session-1-fold-1",
		Path:          "session-1.kv",
		IndexURI:      "mlx://state-ramp/fold/1/folded/index",
		TokenCount:    206,
		PayloadOffset: 1234,
		PayloadBytes:  80511040,
	}
}

func TestNew(t *testing.T) {
	idx := New()
	if idx.Kind != DefaultKind {
		t.Errorf("New().Kind = %q, want %q", idx.Kind, DefaultKind)
	}
	if idx.Len() != 0 {
		t.Errorf("New().Len() = %d, want 0", idx.Len())
	}
}

func TestIndexAppend(t *testing.T) {
	idx := New()
	if err := idx.Append(sampleEntry()); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if idx.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", idx.Len())
	}

	got, ok := idx.Lookup("mlx://state-ramp/fold/1/folded/index")
	if !ok {
		t.Fatal("Lookup() did not find the entry just appended")
	}
	if got != sampleEntry() {
		t.Errorf("Lookup() = %+v, want %+v", got, sampleEntry())
	}
}

func TestIndexAppendRejects(t *testing.T) {
	cases := []struct {
		name  string
		entry Entry
		want  error
	}{
		{"NoIndexURI", Entry{Path: "a.kv"}, ErrMissingIndexURI},
		{"NoPath", Entry{IndexURI: "mlx://a"}, ErrMissingPath},
		{"NegativeOffset", Entry{IndexURI: "mlx://a", Path: "a.kv", PayloadOffset: -1}, ErrNegativeRange},
		{"NegativeBytes", Entry{IndexURI: "mlx://a", Path: "a.kv", PayloadBytes: -1}, ErrNegativeRange},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			idx := New()
			err := idx.Append(tc.entry)
			if !errors.Is(err, tc.want) {
				t.Errorf("Append() error = %v, want %v", err, tc.want)
			}
			if idx.Len() != 0 {
				t.Errorf("a rejected entry was stored anyway")
			}
		})
	}

	t.Run("Duplicate", func(t *testing.T) {
		idx := New()
		if err := idx.Append(sampleEntry()); err != nil {
			t.Fatalf("Append() error = %v", err)
		}
		clash := sampleEntry()
		clash.ID = "a different segment on the same uri"
		if err := idx.Append(clash); !errors.Is(err, ErrDuplicateIndexURI) {
			t.Errorf("Append(duplicate) error = %v, want %v", err, ErrDuplicateIndexURI)
		}
		if idx.Len() != 1 {
			t.Errorf("Len() = %d, want 1 — the duplicate was stored", idx.Len())
		}
	})
}

func TestIndexSet(t *testing.T) {
	idx := New()
	if err := idx.Append(sampleEntry()); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if err := idx.Append(Entry{ID: "second", Path: "session-2.kv", IndexURI: "mlx://second"}); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	// Re-packing a segment replaces it in place rather than appending.
	updated := sampleEntry()
	updated.PayloadOffset = 4321
	updated.PayloadBytes = 999
	if err := idx.Set(updated); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if idx.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", idx.Len())
	}
	if idx.States[0].PayloadOffset != 4321 {
		t.Errorf("Set() did not replace in place: %+v", idx.States[0])
	}

	// A new index URI is appended.
	if err := idx.Set(Entry{ID: "third", Path: "session-3.kv", IndexURI: "mlx://third"}); err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if idx.Len() != 3 {
		t.Errorf("Len() = %d, want 3", idx.Len())
	}
	if err := idx.Set(Entry{Path: "no-uri.kv"}); !errors.Is(err, ErrMissingIndexURI) {
		t.Errorf("Set(invalid) error = %v, want %v", err, ErrMissingIndexURI)
	}
}

func TestIndexRemove(t *testing.T) {
	idx := New()
	if err := idx.Append(sampleEntry()); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if !idx.Remove(sampleEntry().IndexURI) {
		t.Error("Remove() reported nothing removed")
	}
	if idx.Len() != 0 {
		t.Errorf("Len() = %d, want 0", idx.Len())
	}
	if idx.Remove("mlx://never-there") {
		t.Error("Remove() reported a removal that did not happen")
	}
}

// TestIndexManyFiles covers the "one .kv or many" requirement: entries may
// share a container file or each name their own.
func TestIndexManyFiles(t *testing.T) {
	idx := New()
	entries := []Entry{
		{ID: "fold-1", Path: "session-1.kv", IndexURI: "mlx://fold/1", PayloadOffset: 100, PayloadBytes: 10},
		{ID: "fold-2", Path: "session-1.kv", IndexURI: "mlx://fold/2", PayloadOffset: 110, PayloadBytes: 20},
		{ID: "fold-3", Path: "session-2.kv", IndexURI: "mlx://fold/3", PayloadOffset: 100, PayloadBytes: 30},
	}
	for _, e := range entries {
		if err := idx.Append(e); err != nil {
			t.Fatalf("Append(%s) error = %v", e.ID, err)
		}
	}

	shared := idx.ByPath("session-1.kv")
	if len(shared) != 2 || shared[0].ID != "fold-1" || shared[1].ID != "fold-2" {
		t.Errorf("ByPath(session-1.kv) = %+v, want fold-1 and fold-2 in order", shared)
	}
	if own := idx.ByPath("session-2.kv"); len(own) != 1 || own[0].ID != "fold-3" {
		t.Errorf("ByPath(session-2.kv) = %+v, want just fold-3", own)
	}
	if none := idx.ByPath("session-9.kv"); len(none) != 0 {
		t.Errorf("ByPath(unknown) = %+v, want nothing", none)
	}

	paths := idx.Paths()
	want := []string{"session-1.kv", "session-2.kv"}
	if len(paths) != len(want) {
		t.Fatalf("Paths() = %v, want %v", paths, want)
	}
	for n := range want {
		if paths[n] != want[n] {
			t.Errorf("Paths()[%d] = %q, want %q", n, paths[n], want[n])
		}
	}
}

// TestEncodeShape pins the on-disk field names. They are the interop contract
// with whatever writes the containers, so a rename has to be deliberate.
func TestEncodeShape(t *testing.T) {
	idx := New()
	if err := idx.Append(sampleEntry()); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	var buf bytes.Buffer
	if err := idx.Encode(&buf); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &raw); err != nil {
		t.Fatalf("encoded index is not valid JSON: %v", err)
	}
	if raw["kind"] != DefaultKind {
		t.Errorf("kind = %v, want %q", raw["kind"], DefaultKind)
	}
	states, ok := raw["states"].([]interface{})
	if !ok || len(states) != 1 {
		t.Fatalf("states = %v, want one entry", raw["states"])
	}
	entry, ok := states[0].(map[string]interface{})
	if !ok {
		t.Fatalf("states[0] = %v, want an object", states[0])
	}
	for key, want := range map[string]interface{}{
		"id":             "session-1-fold-1",
		"path":           "session-1.kv",
		"index_uri":      "mlx://state-ramp/fold/1/folded/index",
		"token_count":    float64(206),
		"payload_offset": float64(1234),
		"payload_bytes":  float64(80511040),
	} {
		if got := entry[key]; got != want {
			t.Errorf("states[0][%q] = %v, want %v", key, got, want)
		}
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	idx := New()
	for _, e := range []Entry{
		sampleEntry(),
		{ID: "fold-2", Path: "session-2.kv", IndexURI: "mlx://fold/2", PayloadOffset: 28, PayloadBytes: 65 << 20},
	} {
		if err := idx.Append(e); err != nil {
			t.Fatalf("Append() error = %v", err)
		}
	}

	var buf bytes.Buffer
	if err := idx.Encode(&buf); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	got, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got.Kind != idx.Kind || got.Len() != idx.Len() {
		t.Fatalf("Decode() = kind %q, %d entries; want %q, %d", got.Kind, got.Len(), idx.Kind, idx.Len())
	}
	for n := range idx.States {
		if got.States[n] != idx.States[n] {
			t.Errorf("entry %d round-tripped as %+v, want %+v", n, got.States[n], idx.States[n])
		}
	}
}

func TestDecodeRejects(t *testing.T) {
	cases := []struct {
		name string
		json string
		want error
	}{
		{"MissingIndexURI", `{"kind":"x","states":[{"id":"a","path":"a.kv"}]}`, ErrMissingIndexURI},
		{"MissingPath", `{"kind":"x","states":[{"id":"a","index_uri":"mlx://a"}]}`, ErrMissingPath},
		{"NegativeRange", `{"kind":"x","states":[{"path":"a.kv","index_uri":"mlx://a","payload_bytes":-1}]}`, ErrNegativeRange},
		{"Duplicate", `{"kind":"x","states":[{"path":"a.kv","index_uri":"mlx://a"},{"path":"b.kv","index_uri":"mlx://a"}]}`, ErrDuplicateIndexURI},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Decode(strings.NewReader(tc.json)); !errors.Is(err, tc.want) {
				t.Errorf("Decode() error = %v, want %v", err, tc.want)
			}
		})
	}

	t.Run("MalformedJSON", func(t *testing.T) {
		if _, err := Decode(strings.NewReader("{not json")); err == nil {
			t.Error("Decode(malformed) returned no error")
		}
	})
	t.Run("NilReader", func(t *testing.T) {
		if _, err := Decode(nil); err == nil {
			t.Error("Decode(nil) returned no error")
		}
	})
}

// TestIndexNeverReadsThePayload is the point of the package: an index records
// byte ranges, it does not open containers. Every operation works on entries
// naming files that do not exist — and, when one does exist, its bytes are
// never touched.
func TestIndexNeverReadsThePayload(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "never-written.kv")

	idx := New()
	if err := idx.Append(Entry{
		ID:            "absent",
		Path:          missing,
		IndexURI:      "mlx://absent",
		PayloadOffset: 28,
		PayloadBytes:  1 << 40, // a terabyte that was never written anywhere
	}); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	var buf bytes.Buffer
	if err := idx.Encode(&buf); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	got, err := Decode(&buf)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	entry, ok := got.Lookup("mlx://absent")
	if !ok {
		t.Fatal("Lookup() lost an entry whose file does not exist")
	}
	if entry.PayloadBytes != 1<<40 {
		t.Errorf("PayloadBytes = %d, want %d", entry.PayloadBytes, int64(1<<40))
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Errorf("the index materialised %s", missing)
	}

	// The same holds for a file that does exist: the recorded range is taken
	// at its word, never verified against the bytes on disk.
	present := filepath.Join(dir, "present.kv")
	if err := os.WriteFile(present, []byte("KVST"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := idx.Append(Entry{
		ID:            "present",
		Path:          present,
		IndexURI:      "mlx://present",
		PayloadOffset: 28,
		PayloadBytes:  80511040,
	}); err != nil {
		t.Fatalf("Append() error = %v", err)
	}
	if e, ok := idx.Lookup("mlx://present"); !ok || e.PayloadBytes != 80511040 {
		t.Errorf("Lookup() = %+v, %v — the range was second-guessed against the file", e, ok)
	}
}
