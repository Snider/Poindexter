// CommonJS loader placeholder for @snider/poindexter-wasm
// This package is intended for browser bundlers (Angular/webpack/Vite) using ESM.
// If you are in a CommonJS environment, please switch to ESM import:
//   import { init } from '@snider/poindexter-wasm';
// Or configure your bundler to use the ESM entry.

module.exports = {
  init: function () {
    throw new Error("@snider/poindexter-wasm: CommonJS is not supported; use ESM import instead.");
  }
};
