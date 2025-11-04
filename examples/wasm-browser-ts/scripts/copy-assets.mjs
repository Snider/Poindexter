// Copies WASM artifacts and loader into the public/ folder before Vite dev/build.
// Run as an npm script (predev) from this example directory.
import { cp, mkdir } from 'node:fs/promises';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

async function main() {
  const root = resolve(__dirname, '../../..');
  const exampleDir = resolve(__dirname, '..');
  const publicDir = resolve(exampleDir, 'public');

  await mkdir(publicDir, { recursive: true });

  const sources = [
    // WASM artifacts built by `make wasm-build`
    resolve(root, 'dist/poindexter.wasm'),
    resolve(root, 'dist/wasm_exec.js'),
    // ESM loader shipped with the repo's npm folder
    resolve(root, 'npm/poindexter-wasm/loader.js'),
  ];

  const targets = [
    resolve(publicDir, 'poindexter.wasm'),
    resolve(publicDir, 'wasm_exec.js'),
    resolve(publicDir, 'loader.js'),
  ];

  for (let i = 0; i < sources.length; i++) {
    await cp(sources[i], targets[i]);
    console.log(`Copied ${sources[i]} -> ${targets[i]}`);
  }
}

main().catch((err) => {
  console.error('copy-assets failed:', err);
  process.exit(1);
});
