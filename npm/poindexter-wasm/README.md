# @snider/poindexter-wasm

WebAssembly build of the Poindexter KD-Tree library for browsers. Designed to be consumed from Angular, React, or any ESM-capable bundler.

Status: experimental preview. API surface can evolve.

## Install

Until published to npm, you can use a local file/path install:

```bash
# From the repo root where this folder exists
npm pack ./npm/poindexter-wasm
# Produces a tarball like snider-poindexter-wasm-0.0.0-development.tgz
# In your Angular project:
npm install ../Poindexter/snider-poindexter-wasm-0.0.0-development.tgz
```

Once published:

```bash
npm install @snider/poindexter-wasm
```

## Usage (Angular/ESM)

```ts
// app.module.ts or a dedicated provider file
import { init } from '@snider/poindexter-wasm';

async function bootstrapPoindexter() {
  const px = await init();
  console.log(await px.version());

  const tree = await px.newTree(2);
  await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
  await tree.insert({ id: 'b', coords: [1, 1], value: 'B' });

  const nearest = await tree.nearest([0.2, 0.1]);
  console.log('nearest:', nearest);

  return { px, tree };
}

// Call bootstrapPoindexter() during app initialization
```

If your bundler cannot resolve asset URLs from `import.meta.url`, pass explicit URLs:

```ts
const px = await init({
  wasmURL: '/assets/poindexter/poindexter.wasm',
  wasmExecURL: '/assets/poindexter/wasm_exec.js',
});
```

To host the assets, copy `node_modules/@snider/poindexter-wasm/dist/*` into your app's public/assets folder during build (e.g., with Angular `assets` config in `angular.json`).

## API

- `version(): Promise<string>` – Poindexter library version.
- `hello(name?: string): Promise<string>` – simple sanity check.
- `newTree(dim: number): Promise<Tree>` – create a new KD-Tree with given dimension.

Tree methods:
- `dim(): Promise<number>`
- `len(): Promise<number>`
- `insert(point: {id: string, coords: number[], value?: string}): Promise<boolean>`
- `deleteByID(id: string): Promise<boolean>`
- `nearest(query: number[]): Promise<{point, dist, found}>`
- `kNearest(query: number[], k: number): Promise<{points, dists}>`
- `radius(query: number[], r: number): Promise<{points, dists}>`
- `exportJSON(): Promise<string>` – minimal metadata export for now.

## Notes

- Values are strings in this WASM build for simplicity across the boundary.
- This package ships `dist/poindexter.wasm` and Go's `wasm_exec.js`. The loader adds required shims at runtime.
- Requires a modern browser with WebAssembly support.
