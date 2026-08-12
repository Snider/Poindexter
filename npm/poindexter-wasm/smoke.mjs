// Minimal Node smoke test for the WASM loader.
// Assumes npm-pack has prepared npm/poindexter-wasm with loader and dist assets.

import { init } from './loader.js';

(async function () {
  let px;
  try {
    px = await init({
      wasmURL: new URL('./dist/poindexter.wasm', import.meta.url).toString(),
      wasmExecURL: new URL('./dist/wasm_exec.js', import.meta.url).toString(),
    });
    const ver = await px.version();
    if (!ver || typeof ver !== 'string') throw new Error('version not string');

    const tree = await px.newTree(2);
    await tree.insert({ id: 'a', coords: [0, 0], value: 'A' });
    await tree.insert({ id: 'b', coords: [1, 0], value: 'B' });
    const nn = await tree.nearest([0.9, 0.1]);
    if (!nn || !nn.point || !nn.point.id) throw new Error('nearest failed');
    console.log('WASM smoke ok:', ver, 'nearest.id=', nn.point.id);
  } catch (err) {
    console.error('WASM smoke failed:', err);
    process.exit(1);
  }

  // Test error handling
  try {
    await px.newTree(0);
    console.error('Expected error from newTree(0) but got none');
    process.exit(1);
  } catch (err) {
    if (err.code === 'bad_request' && err.message.includes('dimension')) {
      console.log('WASM error handling ok:', err.code);
    } else {
      console.error('WASM smoke failed: unexpected error format:', err);
      process.exit(1);
    }
  }
})();
