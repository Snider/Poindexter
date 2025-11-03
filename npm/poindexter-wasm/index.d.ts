export interface PxPoint {
  id: string;
  coords: number[];
  value?: string;
}

export interface NearestResult {
  point: PxPoint;
  dist: number;
  found: boolean;
}

export interface KNearestResult {
  points: PxPoint[];
  dists: number[];
}

export interface PxTree {
  len(): Promise<number>;
  dim(): Promise<number>;
  insert(point: PxPoint): Promise<boolean>;
  deleteByID(id: string): Promise<boolean>;
  nearest(query: number[]): Promise<NearestResult>;
  kNearest(query: number[], k: number): Promise<KNearestResult>;
  radius(query: number[], r: number): Promise<KNearestResult>;
  exportJSON(): Promise<string>;
}

export interface InitOptions {
  wasmURL?: string;
  wasmExecURL?: string;
  instantiateWasm?: (source: ArrayBuffer, importObject: WebAssembly.Imports) => Promise<WebAssembly.Instance> | WebAssembly.Instance;
}

export interface PxAPI {
  version(): Promise<string>;
  hello(name?: string): Promise<string>;
  newTree(dim: number): Promise<PxTree>;
}

export function init(options?: InitOptions): Promise<PxAPI>;
