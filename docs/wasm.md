# Browser/WebAssembly (WASM)

Poindexter ships a browser build compiled to WebAssembly along with a small JS loader and TypeScript types. This allows you to use the KD‑Tree functionality directly from web apps (Angular, React, Vue, plain ESM, etc.).

## What’s included

- `dist/poindexter.wasm` — the compiled Go WASM module
- `dist/wasm_exec.js` — Go’s runtime shim required to run WASM in the browser
- `npm/poindexter-wasm/loader.js` — ESM loader that instantiates the WASM and exposes a friendly API
- `npm/poindexter-wasm/index.d.ts` — TypeScript typings for the loader and KD‑Tree API

## Quick start

- Build artifacts and copy `wasm_exec.js`:

```bash
make wasm-build
```

- Prepare the npm package folder with `dist/` and docs:

```bash
make npm-pack
```

- Minimal browser ESM usage (serve `dist/` statically):

```html
<script type="module">
  import { init } from '/npm/poindexter-wasm/loader.js';
  const px = await init({
    wasmURL: '/dist/poindexter.wasm',
    wasmExecURL: '/dist/wasm_exec.js',
  });
  const tree = await px.newTree(2);
  await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
  const nn = await tree.nearest([0.1, 0.2]);
  console.log(nn);
</script>
```

## Building locally

```bash
make wasm-build
```

This produces `dist/poindexter.wasm` and copies `wasm_exec.js` into `dist/` from your Go installation. If your environment is non‑standard, you can override the path:

```bash
WASM_EXEC=/custom/path/wasm_exec.js make wasm-build
```

To assemble the npm package folder with the built artifacts:

```bash
make npm-pack
```

This populates `npm/poindexter-wasm/` with `dist/`, licence and readme files. You can then create a tarball for local testing:

```bash
npm pack ./npm/poindexter-wasm
```

## Using in Angular (example)

1) Install the package (use the tarball generated above or a published version):

```bash
npm install <path-to>/snider-poindexter-wasm-0.0.0-development.tgz
# or once published
npm install @snider/poindexter-wasm
```

2) Make the WASM runtime files available as app assets. In `angular.json` under `build.options.assets`:

```json
{
  "glob": "**/*",
  "input": "node_modules/@snider/poindexter-wasm/dist",
  "output": "/assets/poindexter/"
}
```

3) Import and initialize in your code:

```ts
import { init } from '@snider/poindexter-wasm';

const px = await init({
  // If you used the assets mapping above, these defaults should work:
  wasmURL: '/assets/poindexter/poindexter.wasm',
  wasmExecURL: '/assets/poindexter/wasm_exec.js',
});

const tree = await px.newTree(2);
await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
const nearest = await tree.nearest([0.1, 0.2]);
console.log(nearest);
```

## JavaScript API

Top‑level functions returned by `init()`:

- `version(): string`
- `hello(name?: string): string`
- `newTree(dim: number): Promise<Tree>`

Tree methods:

- `dim(): Promise<number>`
- `len(): Promise<number>`
- `insert(p: { id: string; coords: number[]; value?: string }): Promise<void>`
- `deleteByID(id: string): Promise<boolean>`
- `nearest(query: number[]): Promise<{ id: string; coords: number[]; value: string; dist: number } | null>`
- `kNearest(query: number[], k: number): Promise<Array<{ id: string; coords: number[]; value: string; dist: number }>>`
- `radius(query: number[], r: number): Promise<Array<{ id: string; coords: number[]; value: string; dist: number }>>`
- `exportJSON(): Promise<string>`

Notes:
- The WASM bridge currently uses `KDTree[string]` for values to keep the boundary simple. You can encode richer payloads as JSON strings if needed.
- `wasm_exec.js` must be available next to the `.wasm` file (the loader accepts explicit URLs if you place them elsewhere).

## CI artifacts

Our CI builds and uploads the following artifacts on each push/PR:

- `poindexter-wasm-dist` — the `dist/` folder containing `poindexter.wasm` and `wasm_exec.js`
- `npm-poindexter-wasm` — the prepared npm package folder with `dist/` and documentation
- `npm-poindexter-wasm-tarball` — a `.tgz` created via `npm pack` for quick local install/testing

You can download these artifacts from the workflow run summary in GitHub Actions.

## Browser demo (checked into repo)

There is a tiny browser demo you can load locally from this repo:

- Path: `examples/wasm-browser/index.html`
- Prerequisites: run `make wasm-build` so `dist/poindexter.wasm` and `dist/wasm_exec.js` exist.
- Serve the repo root (so relative paths resolve), for example:

```bash
python3 -m http.server -b 127.0.0.1 8000
```

Then open:

- http://127.0.0.1:8000/examples/wasm-browser/

Open the browser console to see outputs from `nearest`, `kNearest`, and `radius` queries.
