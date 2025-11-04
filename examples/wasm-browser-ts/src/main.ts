// Minimal TypeScript demo that uses the Poindexter WASM ESM loader.
// Precondition: run `make wasm-build` at repo root, then `npm run dev` in this folder.

// We copy the loader and wasm artifacts to /public via scripts/copy-assets.mjs before dev starts.
// @ts-ignore
import { init } from '/loader.js';

async function run() {
  const px = await init({
    wasmURL: '/poindexter.wasm',
    wasmExecURL: '/wasm_exec.js',
  });

  console.log('Poindexter (WASM) version:', await px.version());

  const tree = await px.newTree(2);
  await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
  await tree.insert({ id: 'b', coords: [1, 0], value: 'B' });
  await tree.insert({ id: 'c', coords: [0, 1], value: 'C' });

  const nn = await tree.nearest([0.9, 0.1]);
  console.log('Nearest [0.9,0.1]:', nn);

  const kn = await tree.kNearest([0.9, 0.9], 2);
  console.log('kNN k=2 [0.9,0.9]:', kn);

  const rad = await tree.radius([0, 0], 1.1);
  console.log('Radius r=1.1 [0,0]:', rad);
}

run().catch((err) => {
  console.error('WASM demo error:', err);
});
