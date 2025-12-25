// ESM loader for Poindexter WASM
// Usage:
//   import { init } from '@snider/poindexter-wasm';
//   const px = await init();
//   const tree = await px.newTree(2);
//   await tree.insert({ id: 'a', coords: [0,0], value: 'A' });
//   const res = await tree.nearest([0.1, 0.2]);

async function loadScriptOnce(src) {
  return new Promise((resolve, reject) => {
    // If already present, resolve immediately
    if (document.querySelector(`script[src="${src}"]`)) return resolve();
    const s = document.createElement('script');
    s.src = src;
    s.onload = () => resolve();
    s.onerror = (e) => reject(new Error(`Failed to load ${src}`));
    document.head.appendChild(s);
  });
}

async function ensureWasmExec(url) {
  if (typeof window !== 'undefined' && typeof window.Go === 'function') return;
  await loadScriptOnce(url);
  if (typeof window === 'undefined' || typeof window.Go !== 'function') {
    throw new Error('wasm_exec.js did not define window.Go');
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
  // Core operations
  async len() { return call('pxTreeLen', this.treeId); }
  async dim() { return call('pxTreeDim', this.treeId); }
  async insert(point) { return call('pxInsert', this.treeId, point); }
  async deleteByID(id) { return call('pxDeleteByID', this.treeId, id); }
  async nearest(query) { return call('pxNearest', this.treeId, query); }
  async kNearest(query, k) { return call('pxKNearest', this.treeId, query, k); }
  async radius(query, r) { return call('pxRadius', this.treeId, query, r); }
  async exportJSON() { return call('pxExportJSON', this.treeId); }
  // Analytics operations
  async getAnalytics() { return call('pxGetAnalytics', this.treeId); }
  async getPeerStats() { return call('pxGetPeerStats', this.treeId); }
  async getTopPeers(n) { return call('pxGetTopPeers', this.treeId, n); }
  async getAxisDistributions(axisNames) { return call('pxGetAxisDistributions', this.treeId, axisNames); }
  async resetAnalytics() { return call('pxResetAnalytics', this.treeId); }
}

export async function init(options = {}) {
  const {
    wasmURL = new URL('./dist/poindexter.wasm', import.meta.url).toString(),
    wasmExecURL = new URL('./dist/wasm_exec.js', import.meta.url).toString(),
    instantiateWasm // optional custom instantiator: (source, importObject) => WebAssembly.Instance
  } = options;

  await ensureWasmExec(wasmExecURL);
  const go = new window.Go();

  let result;
  if (instantiateWasm) {
    const source = await fetch(wasmURL).then(r => r.arrayBuffer());
    const inst = await instantiateWasm(source, go.importObject);
    result = { instance: inst };
  } else if (WebAssembly.instantiateStreaming) {
    result = await WebAssembly.instantiateStreaming(fetch(wasmURL), go.importObject);
  } else {
    const resp = await fetch(wasmURL);
    const bytes = await resp.arrayBuffer();
    result = await WebAssembly.instantiate(bytes, go.importObject);
  }

  // Run the Go program (it registers globals like pxNewTree, etc.)
  // Do not await: the Go WASM main may block (e.g., via select{}), so awaiting never resolves.
  go.run(result.instance);

  const api = {
    // Core functions
    version: async () => call('pxVersion'),
    hello: async (name) => call('pxHello', name ?? ''),
    newTree: async (dim) => {
      const info = call('pxNewTree', dim);
      return new PxTree(info.treeId);
    },
    // Statistics utilities
    computeDistributionStats: async (distances) => call('pxComputeDistributionStats', distances),
    // NAT routing / peer quality functions
    computePeerQualityScore: async (metrics, weights) => call('pxComputePeerQualityScore', metrics, weights),
    computeTrustScore: async (metrics) => call('pxComputeTrustScore', metrics),
    getDefaultQualityWeights: async () => call('pxGetDefaultQualityWeights'),
    getDefaultPeerFeatureRanges: async () => call('pxGetDefaultPeerFeatureRanges'),
    normalizePeerFeatures: async (features, ranges) => call('pxNormalizePeerFeatures', features, ranges),
    weightedPeerFeatures: async (normalized, weights) => call('pxWeightedPeerFeatures', normalized, weights)
  };

  return api;
}
