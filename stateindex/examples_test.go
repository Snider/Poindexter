package stateindex_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Snider/Poindexter/stateindex"
)

func ExampleIndex_Append() {
	idx := stateindex.New()
	if err := idx.Append(stateindex.Entry{
		ID:            "session-1-fold-1",
		Path:          "session-1.kv",
		IndexURI:      "mlx://state-ramp/fold/1/folded/index",
		TokenCount:    206,
		PayloadOffset: 1234,
		PayloadBytes:  80511040,
	}); err != nil {
		log.Fatalf("Append failed: %v", err)
	}

	entry, ok := idx.Lookup("mlx://state-ramp/fold/1/folded/index")
	if !ok {
		log.Fatal("entry not found")
	}
	// The byte range came from the index, not from opening session-1.kv.
	fmt.Printf("%s: %s [%d..%d)\n", entry.ID, entry.Path, entry.PayloadOffset, entry.PayloadOffset+entry.PayloadBytes)
	// Output:
	// session-1-fold-1: session-1.kv [1234..80512274)
}

func ExampleIndex_Paths() {
	idx := stateindex.New()
	for _, e := range []stateindex.Entry{
		{ID: "fold-1", Path: "session-1.kv", IndexURI: "mlx://fold/1", PayloadOffset: 28, PayloadBytes: 4096},
		{ID: "fold-2", Path: "session-1.kv", IndexURI: "mlx://fold/2", PayloadOffset: 4124, PayloadBytes: 8192},
		{ID: "fold-3", Path: "session-2.kv", IndexURI: "mlx://fold/3", PayloadOffset: 28, PayloadBytes: 2048},
	} {
		if err := idx.Append(e); err != nil {
			log.Fatalf("Append failed: %v", err)
		}
	}

	// One index can span several container files, or share one between many
	// segments — Paths reports the files a consumer actually has to open.
	fmt.Println(idx.Paths())
	fmt.Println(len(idx.ByPath("session-1.kv")))
	// Output:
	// [session-1.kv session-2.kv]
	// 2
}

func ExampleDecode() {
	sidecar := `{
  "kind": "lthn/state-index",
  "states": [
    {
      "id": "session-1-fold-1",
      "path": "session-1.kv",
      "index_uri": "mlx://state-ramp/fold/1/folded/index",
      "token_count": 206,
      "payload_offset": 1234,
      "payload_bytes": 80511040
    }
  ]
}`

	idx, err := stateindex.Decode(strings.NewReader(sidecar))
	if err != nil {
		log.Fatalf("Decode failed: %v", err)
	}
	if err := idx.Encode(os.Stdout); err != nil {
		log.Fatalf("Encode failed: %v", err)
	}
	// Output:
	// {
	//   "kind": "lthn/state-index",
	//   "states": [
	//     {
	//       "id": "session-1-fold-1",
	//       "path": "session-1.kv",
	//       "index_uri": "mlx://state-ramp/fold/1/folded/index",
	//       "token_count": 206,
	//       "payload_offset": 1234,
	//       "payload_bytes": 80511040
	//     }
	//   ]
	// }
}
