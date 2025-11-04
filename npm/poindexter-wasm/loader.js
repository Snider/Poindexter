// ESM loader for Poindexter WASM
// Usage:
//   import { init } from '@snider/poindexter-wasm';
//   const px = await init();
//   const tree = await px.newTree(2);
//   await tree.insert({ id: 'a', coords: [0,0], value: 'A' });
//   const res = await tree.nearest([0.1, 0.2]);

const isNode = typeof process !== 'undefined' && process.versions && process.versions.node;

async function loadScriptOnce(src) {
  // Browser-only
  return new Promise((resolve, reject) => {
    if (typeof document === 'undefined') return reject(new Error('loadScriptOnce requires a browser environment'));
    if (document.querySelector(`script[src="${src}"]`)) return resolve();
    const s = document.createElement('script');
    s.src = src;
    s.onload = () => resolve();
    s.onerror = (e) => reject(new Error(`Failed to load ${src}`));
    document.head.appendChild(s);
  });
}

async function ensureWasmExec(url) {
  if (globalThis.Go) return;

  if (isNode) {
    const { fileURLToPath } = await import('url');
    const fs = await import('fs/promises');
    const vm = await import('vm');
    const wasmExecPath = fileURLToPath(url);
    const wasmExecCode = await fs.readFile(wasmExecPath, 'utf8');
    vm.runInThisContext(wasmExecCode, { filename: wasmExecPath });
  } else if (typeof window !== 'undefined') {
    await loadScriptOnce(url);
  } else {
    throw new Error('Unsupported environment: not Node.js or a browser');
  }

  if (typeof globalThis.Go !== 'function') {
    throw new Error('wasm_exec.js did not define globalThis.Go');
  }
}

function unwrap(result) {
  if (!result || typeof result !== 'object') throw new Error('bad result');
  if (result.ok) return result.data;
  throw new Error(result.error || 'unknown error');
}

function call(name, ...args) {
  const fn = globalThis[name];
  if (typeof fn !== 'function') throw new Error(`WASM function ${name} not found`);
  return unwrap(fn(...args));
}

class PxTree {
  constructor(treeId) { this.treeId = treeId; }
  async len() { return call('pxTreeLen', this.treeId); }
  async dim() { return call('pxTreeDim', this.treeId); }
  async insert(point) { return call('pxInsert', this.treeId, point); }
  async deleteByID(id) { return call('pxDeleteByID', this.treeId, id); }
  async nearest(query) {
    const res = await call('pxNearest', this.treeId, query);
    if (!res.found) return { point: null, dist: 0, found: false };
    return {
      point: { id: res.id, coords: res.coords, value: res.value },
      dist: res.dist,
      found: true,
    };
  }
  async kNearest(query, k) { return call('pxKNearest', this.treeId, query, k); }
  async radius(query, r) { return call('pxRadius', this.treeId, query, r); }
  async exportJSON() { return call('pxExportJSON', this.treeId); }
}

export async function init(options = {}) {
  const {
    wasmURL = new URL('./dist/poindexter.wasm', import.meta.url).toString(),
    wasmExecURL = new URL('./dist/wasm_exec.js', import.meta.url).toString(),
    instantiateWasm, // optional custom instantiator
    fetch: customFetch,
  } = options;

  await ensureWasmExec(wasmExecURL);
  const go = new globalThis.Go();

  const fetchFn = customFetch || (isNode ? (async (url) => {
    const { fileURLToPath } = await import('url');
    const fs = await import('fs/promises');
    const path = fileURLToPath(url);
    const bytes = await fs.readFile(path);
    return new Response(bytes, { 'Content-Type': 'application/wasm' });
  }) : fetch);


  let result;
  if (instantiateWasm) {
    const source = await fetchFn(wasmURL).then(r => r.arrayBuffer());
    const inst = await instantiateWasm(source, go.importObject);
    result = { instance: inst };
  } else if (WebAssembly.instantiateStreaming && !isNode) {
    result = await WebAssembly.instantiateStreaming(fetchFn(wasmURL), go.importObject);
  } else {
    const resp = await fetchFn(wasmURL);
    const bytes = await resp.arrayBuffer();
    result = await WebAssembly.instantiate(bytes, go.importObject);
  }

  // Run the Go program (it registers globals like pxNewTree, etc.)
  // Do not await: the Go WASM main may block (e.g., via select{}), so awaiting never resolves.
  go.run(result.instance);

  const api = {
    version: async () => call('pxVersion'),
    hello: async (name) => call('pxHello', name ?? ''),
    newTree: async (dim) => {
      const info = call('pxNewTree', dim);
      return new PxTree(info.treeId);
    },
  };

  return api;
}
