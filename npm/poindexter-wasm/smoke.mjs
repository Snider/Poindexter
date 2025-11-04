// Minimal Node smoke test for the WASM loader.
// Assumes npm-pack has prepared npm/poindexter-wasm with loader and dist assets.

import { init } from './loader.js';

(async function () {
  try {
    const px = await init({
      // In CI, dist/ is placed at repo root via make wasm-build && make npm-pack
      wasmURL: new URL('./dist/poindexter.wasm', import.meta.url).pathname,
      wasmExecURL: new URL('./dist/wasm_exec.js', import.meta.url).pathname,
    });
    const ver = await px.version();
    if (!ver || typeof ver !== 'string') throw new Error('version not string');

    const tree = await px.newTree(2);
    await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
    await tree.insert({ id: 'b', coords: [1, 0], value: 'B' });
    const nn = await tree.nearest([0.9, 0.1]);
    if (!nn || !nn.id) throw new Error('nearest failed');
    console.log('WASM smoke ok:', ver, 'nearest.id=', nn.id);
  } catch (err) {
    console.error('WASM smoke failed:', err);
    process.exit(1);
  }
})();
