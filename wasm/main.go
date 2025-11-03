//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	pd "github.com/Snider/Poindexter"
)

// Simple registry for KDTree instances created from JS.
// We keep values as string for simplicity across the WASM boundary.
var (
	treeRegistry = map[int]*pd.KDTree[string]{}
	nextTreeID   = 1
)

func export(name string, fn func(this js.Value, args []js.Value) (any, error)) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		res, err := fn(this, args)
		if err != nil {
			return map[string]any{"ok": false, "error": err.Error()}
		}
		return map[string]any{"ok": true, "data": res}
	}))
}

func getInt(v js.Value, idx int) (int, error) {
	if len := v.Length(); len > idx {
		return v.Index(idx).Int(), nil
	}
	return 0, errors.New("missing integer argument")
}

func getFloatSlice(arg js.Value) ([]float64, error) {
	if arg.IsUndefined() || arg.IsNull() {
		return nil, errors.New("coords/query is undefined or null")
	}
	ln := arg.Length()
	res := make([]float64, ln)
	for i := 0; i < ln; i++ {
		res[i] = arg.Index(i).Float()
	}
	return res, nil
}

func version(_ js.Value, _ []js.Value) (any, error) {
	return pd.Version(), nil
}

func hello(_ js.Value, args []js.Value) (any, error) {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	return pd.Hello(name), nil
}

func newTree(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("newTree(dim) requires dim")
	}
	dim := args[0].Int()
	if dim <= 0 {
		return nil, pd.ErrZeroDim
	}
	t, err := pd.NewKDTreeFromDim[string](dim)
	if err != nil {
		return nil, err
	}
	id := nextTreeID
	nextTreeID++
	treeRegistry[id] = t
	return map[string]any{"treeId": id, "dim": dim}, nil
}

func treeLen(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("len(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.Len(), nil
}

func treeDim(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("dim(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.Dim(), nil
}

func insert(_ js.Value, args []js.Value) (any, error) {
	// insert(treeId, {id: string, coords: number[], value?: string})
	if len(args) < 2 {
		return nil, errors.New("insert(treeId, point)")
	}
	id := args[0].Int()
	pt := args[1]
	pid := pt.Get("id").String()
	coords, err := getFloatSlice(pt.Get("coords"))
	if err != nil {
		return nil, err
	}
	val := pt.Get("value").String()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	okIns := t.Insert(pd.KDPoint[string]{ID: pid, Coords: coords, Value: val})
	return okIns, nil
}

func deleteByID(_ js.Value, args []js.Value) (any, error) {
	// deleteByID(treeId, id)
	if len(args) < 2 {
		return nil, errors.New("deleteByID(treeId, id)")
	}
	id := args[0].Int()
	pid := args[1].String()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.DeleteByID(pid), nil
}

func nearest(_ js.Value, args []js.Value) (any, error) {
	// nearest(treeId, query:number[]) -> {point, dist, found}
	if len(args) < 2 {
		return nil, errors.New("nearest(treeId, query)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	p, d, found := t.Nearest(query)
	out := map[string]any{
		"point": map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value},
		"dist":  d,
		"found": found,
	}
	return out, nil
}

func kNearest(_ js.Value, args []js.Value) (any, error) {
	// kNearest(treeId, query:number[], k:int) -> {points:[...], dists:[...]}
	if len(args) < 3 {
		return nil, errors.New("kNearest(treeId, query, k)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	k := args[2].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	pts, dists := t.KNearest(query, k)
	jsPts := make([]any, len(pts))
	for i, p := range pts {
		jsPts[i] = map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value}
	}
	return map[string]any{"points": jsPts, "dists": dists}, nil
}

func radius(_ js.Value, args []js.Value) (any, error) {
	// radius(treeId, query:number[], r:number) -> {points:[...], dists:[...]}
	if len(args) < 3 {
		return nil, errors.New("radius(treeId, query, r)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	r := args[2].Float()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	pts, dists := t.Radius(query, r)
	jsPts := make([]any, len(pts))
	for i, p := range pts {
		jsPts[i] = map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value}
	}
	return map[string]any{"points": jsPts, "dists": dists}, nil
}

func exportJSON(_ js.Value, args []js.Value) (any, error) {
	// exportJSON(treeId) -> string (all points)
	if len(args) < 1 {
		return nil, errors.New("exportJSON(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	// naive export: ask for all points by radius from origin with large r; or keep
	// internal slice? KDTree doesn't expose iteration, so skip heavy export here.
	// Return metrics only for now.
	m := map[string]any{"dim": t.Dim(), "len": t.Len()}
	b, _ := json.Marshal(m)
	return string(b), nil
}

func main() {
	// Export core API
	export("pxVersion", version)
	export("pxHello", hello)
	export("pxNewTree", newTree)
	export("pxTreeLen", treeLen)
	export("pxTreeDim", treeDim)
	export("pxInsert", insert)
	export("pxDeleteByID", deleteByID)
	export("pxNearest", nearest)
	export("pxKNearest", kNearest)
	export("pxRadius", radius)
	export("pxExportJSON", exportJSON)

	// Keep running
	select {}
}
