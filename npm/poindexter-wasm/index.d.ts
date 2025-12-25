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

// ============================================================================
// Analytics Types
// ============================================================================

/** Tree operation analytics snapshot */
export interface TreeAnalytics {
  queryCount: number;
  insertCount: number;
  deleteCount: number;
  avgQueryTimeNs: number;
  minQueryTimeNs: number;
  maxQueryTimeNs: number;
  lastQueryTimeNs: number;
  lastQueryAt: number; // Unix milliseconds
  createdAt: number; // Unix milliseconds
  backendRebuildCount: number;
  lastRebuiltAt: number; // Unix milliseconds
}

/** Per-peer selection statistics */
export interface PeerStats {
  peerId: string;
  selectionCount: number;
  avgDistance: number;
  lastSelectedAt: number; // Unix milliseconds
}

/** Statistical distribution analysis */
export interface DistributionStats {
  count: number;
  min: number;
  max: number;
  mean: number;
  median: number;
  stdDev: number;
  p25: number;
  p75: number;
  p90: number;
  p99: number;
  variance: number;
  skewness: number;
  sampleSize?: number;
  computedAt?: number; // Unix milliseconds
}

/** Per-axis distribution in the KD-Tree */
export interface AxisDistribution {
  axis: number;
  name: string;
  stats: DistributionStats;
}

// ============================================================================
// NAT Routing Types
// ============================================================================

/** NAT type classification for routing decisions */
export type NATTypeClassification =
  | 'open'
  | 'full_cone'
  | 'restricted_cone'
  | 'port_restricted'
  | 'symmetric'
  | 'symmetric_udp'
  | 'cgnat'
  | 'firewalled'
  | 'relay_required'
  | 'unknown';

/** Network metrics for NAT routing decisions */
export interface NATRoutingMetrics {
  connectivityScore: number; // 0-1: higher = better reachability
  symmetryScore: number; // 0-1: higher = more symmetric NAT
  relayProbability: number; // 0-1: likelihood peer needs relay
  directSuccessRate: number; // 0-1: historical direct connection success
  avgRttMs: number; // Average RTT in milliseconds
  jitterMs: number; // RTT variance in milliseconds
  packetLossRate: number; // 0-1: packet loss rate
  bandwidthMbps: number; // Bandwidth estimate in Mbps
  natType: NATTypeClassification;
  lastProbeAt?: number; // Unix milliseconds
}

/** Weights for peer quality scoring */
export interface QualityWeights {
  latency: number;
  jitter: number;
  packetLoss: number;
  bandwidth: number;
  connectivity: number;
  symmetry: number;
  directSuccess: number;
  relayPenalty: number;
  natType: number;
}

/** Trust metrics for peer reputation */
export interface TrustMetrics {
  reputationScore: number; // 0-1: aggregated trust score
  successfulTransactions: number;
  failedTransactions: number;
  ageSeconds: number; // How long this peer has been known
  lastSuccessAt?: number; // Unix milliseconds
  lastFailureAt?: number; // Unix milliseconds
  vouchCount: number; // Peers vouching for this peer
  flagCount: number; // Reports against this peer
  proofOfWork: number; // Computational proof of stake/work
}

/** Axis min/max range for normalization */
export interface AxisRange {
  min: number;
  max: number;
}

/** Feature ranges for peer feature normalization */
export interface FeatureRanges {
  ranges: AxisRange[];
  labels?: string[];
}

/** Standard peer features for KD-Tree based selection */
export interface StandardPeerFeatures {
  latencyMs: number;
  hopCount: number;
  geoDistanceKm: number;
  trustScore: number;
  bandwidthMbps: number;
  packetLossRate: number;
  connectivityPct: number;
  natScore: number;
}

/** Export data with all points */
export interface TreeExport {
  dim: number;
  len: number;
  backend: string;
  points: PxPoint[];
}

// ============================================================================
// Tree Interface
// ============================================================================

export interface PxTree {
  // Core operations
  len(): Promise<number>;
  dim(): Promise<number>;
  insert(point: PxPoint): Promise<boolean>;
  deleteByID(id: string): Promise<boolean>;
  nearest(query: number[]): Promise<NearestResult>;
  kNearest(query: number[], k: number): Promise<KNearestResult>;
  radius(query: number[], r: number): Promise<KNearestResult>;
  exportJSON(): Promise<string>;

  // Analytics operations
  getAnalytics(): Promise<TreeAnalytics>;
  getPeerStats(): Promise<PeerStats[]>;
  getTopPeers(n: number): Promise<PeerStats[]>;
  getAxisDistributions(axisNames?: string[]): Promise<AxisDistribution[]>;
  resetAnalytics(): Promise<boolean>;
}

// ============================================================================
// Init Options
// ============================================================================

export interface InitOptions {
  wasmURL?: string;
  wasmExecURL?: string;
  instantiateWasm?: (source: ArrayBuffer, importObject: WebAssembly.Imports) => Promise<WebAssembly.Instance> | WebAssembly.Instance;
}

// ============================================================================
// Main API
// ============================================================================

export interface PxAPI {
  // Core functions
  version(): Promise<string>;
  hello(name?: string): Promise<string>;
  newTree(dim: number): Promise<PxTree>;

  // Statistics utilities
  computeDistributionStats(distances: number[]): Promise<DistributionStats>;

  // NAT routing / peer quality functions
  computePeerQualityScore(metrics: NATRoutingMetrics, weights?: QualityWeights): Promise<number>;
  computeTrustScore(metrics: TrustMetrics): Promise<number>;
  getDefaultQualityWeights(): Promise<QualityWeights>;
  getDefaultPeerFeatureRanges(): Promise<FeatureRanges>;
  normalizePeerFeatures(features: number[], ranges?: FeatureRanges): Promise<number[]>;
  weightedPeerFeatures(normalized: number[], weights: number[]): Promise<number[]>;
}

export function init(options?: InitOptions): Promise<PxAPI>;
