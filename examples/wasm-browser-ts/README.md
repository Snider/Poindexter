# WASM Browser Example (TypeScript + Vite)

This is a minimal TypeScript example that runs Poindexter’s WebAssembly build in the browser.
It bundles a tiny page with Vite and demonstrates creating a KDTree and running `Nearest`,
`KNearest`, and `Radius` queries.

## Prerequisites
- Go toolchain installed
- Node.js 18+ (tested with Node 20)

## Quick start

1) Build the WASM artifacts at the repo root:

```bash
make wasm-build
```

This creates `dist/poindexter.wasm` and `dist/wasm_exec.js`.

2) From this example directory, install deps and start the dev server (the script copies the required files into `public/` before starting Vite):

```bash
npm install
npm run dev
```

3) Open the URL printed by Vite (usually http://127.0.0.1:5173/). Open the browser console to see outputs.

## What the dev script does
- Copies `../../dist/poindexter.wasm` and `../../dist/wasm_exec.js` into `public/`
- Copies `../../npm/poindexter-wasm/loader.js` into `public/`
- Starts Vite with `public/` as the static root for those assets

The TypeScript code imports the loader from `/loader.js` and initializes with:

```ts
const px = await init({
  wasmURL: '/poindexter.wasm',
  wasmExecURL: '/wasm_exec.js',
});
```

## Notes
- This example is local-only and not built in CI to keep jobs light.
- You can adapt the same structure inside your own web projects; alternatively, install the published npm package when available and serve `dist/` as static assets.
