import { defineConfig } from 'vite';

// Minimal Vite config for the WASM TS example.
// Serves files from project root; our dev script copies required artifacts to public/.
export default defineConfig({
  root: '.',
  server: {
    host: '127.0.0.1',
    port: 5173,
    open: false,
  },
  preview: {
    host: '127.0.0.1',
    port: 5173,
    open: false,
  },
});
