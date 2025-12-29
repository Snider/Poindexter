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

  // =========================================================================
  // Basic KD-Tree operations
  // =========================================================================
  const tree = await px.newTree(2);
  await tree.insert({ id: 'peer-a', coords: [0, 0], value: 'Peer A' });
  await tree.insert({ id: 'peer-b', coords: [1, 0], value: 'Peer B' });
  await tree.insert({ id: 'peer-c', coords: [0, 1], value: 'Peer C' });
  await tree.insert({ id: 'peer-d', coords: [0.5, 0.5], value: 'Peer D' });

  console.log('\n=== Basic Queries ===');
  const nn = await tree.nearest([0.9, 0.1]);
  console.log('Nearest [0.9,0.1]:', nn);

  const kn = await tree.kNearest([0.5, 0.5], 3);
  console.log('kNN k=3 [0.5,0.5]:', kn);

  const rad = await tree.radius([0, 0], 1.1);
  console.log('Radius r=1.1 [0,0]:', rad);

  // =========================================================================
  // Analytics Demo
  // =========================================================================
  console.log('\n=== Tree Analytics ===');

  // Perform more queries to generate analytics
  for (let i = 0; i < 10; i++) {
    await tree.nearest([Math.random(), Math.random()]);
  }
  await tree.kNearest([0.2, 0.8], 2);
  await tree.kNearest([0.7, 0.3], 2);

  // Get tree-level analytics
  const analytics = await tree.getAnalytics();
  console.log('Tree Analytics:', {
    queryCount: analytics.queryCount,
    insertCount: analytics.insertCount,
    avgQueryTimeNs: analytics.avgQueryTimeNs,
    minQueryTimeNs: analytics.minQueryTimeNs,
    maxQueryTimeNs: analytics.maxQueryTimeNs,
  });

  // =========================================================================
  // Peer Selection Analytics
  // =========================================================================
  console.log('\n=== Peer Selection Analytics ===');

  // Get all peer stats
  const peerStats = await tree.getPeerStats();
  console.log('All Peer Stats:', peerStats);

  // Get top 3 most frequently selected peers
  const topPeers = await tree.getTopPeers(3);
  console.log('Top 3 Peers:', topPeers);

  // =========================================================================
  // Axis Distribution Analysis
  // =========================================================================
  console.log('\n=== Axis Distributions ===');

  const axisDists = await tree.getAxisDistributions(['latency', 'hops']);
  console.log('Axis Distributions:', axisDists);

  // =========================================================================
  // NAT Routing / Peer Quality Scoring
  // =========================================================================
  console.log('\n=== NAT Routing & Peer Quality ===');

  // Simulate peer network metrics
  const peerMetrics = {
    connectivityScore: 0.9,
    symmetryScore: 0.8,
    relayProbability: 0.1,
    directSuccessRate: 0.95,
    avgRttMs: 50,
    jitterMs: 10,
    packetLossRate: 0.01,
    bandwidthMbps: 100,
    natType: 'full_cone' as const,
  };

  const qualityScore = await px.computePeerQualityScore(peerMetrics);
  console.log('Peer Quality Score (0-1):', qualityScore.toFixed(3));

  // Get default quality weights
  const defaultWeights = await px.getDefaultQualityWeights();
  console.log('Default Quality Weights:', defaultWeights);

  // =========================================================================
  // Trust Score Calculation
  // =========================================================================
  console.log('\n=== Trust Score ===');

  const trustMetrics = {
    reputationScore: 0.8,
    successfulTransactions: 150,
    failedTransactions: 3,
    ageSeconds: 86400 * 30, // 30 days
    vouchCount: 5,
    flagCount: 0,
    proofOfWork: 0.5,
  };

  const trustScore = await px.computeTrustScore(trustMetrics);
  console.log('Trust Score (0-1):', trustScore.toFixed(3));

  // =========================================================================
  // Distribution Statistics
  // =========================================================================
  console.log('\n=== Distribution Statistics ===');

  // Simulate some distance measurements
  const distances = [0.1, 0.15, 0.2, 0.25, 0.3, 0.35, 0.4, 0.5, 0.8, 1.2];
  const distStats = await px.computeDistributionStats(distances);
  console.log('Distance Distribution Stats:', {
    count: distStats.count,
    min: distStats.min.toFixed(3),
    max: distStats.max.toFixed(3),
    mean: distStats.mean.toFixed(3),
    median: distStats.median.toFixed(3),
    stdDev: distStats.stdDev.toFixed(3),
    p90: distStats.p90.toFixed(3),
  });

  // =========================================================================
  // Feature Normalization for KD-Tree
  // =========================================================================
  console.log('\n=== Feature Normalization ===');

  // Raw peer features: [latency_ms, hops, geo_km, trust_inv, bw_inv, loss, conn_inv, nat_inv]
  const rawFeatures = [100, 5, 500, 0.1, 50, 0.02, 5, 0.1];

  // Get default feature ranges
  const featureRanges = await px.getDefaultPeerFeatureRanges();
  console.log('Feature Labels:', featureRanges.labels);

  // Normalize features
  const normalizedFeatures = await px.normalizePeerFeatures(rawFeatures);
  console.log('Normalized Features:', normalizedFeatures.map((f: number) => f.toFixed(3)));

  // Apply custom weights
  const customWeights = [1.5, 1.0, 0.5, 1.2, 0.8, 2.0, 1.0, 0.7];
  const weightedFeatures = await px.weightedPeerFeatures(normalizedFeatures, customWeights);
  console.log('Weighted Features:', weightedFeatures.map((f: number) => f.toFixed(3)));

  // =========================================================================
  // Analytics Reset
  // =========================================================================
  console.log('\n=== Analytics Reset ===');
  await tree.resetAnalytics();
  const resetAnalytics = await tree.getAnalytics();
  console.log('After Reset - Query Count:', resetAnalytics.queryCount);

  console.log('\n=== Demo Complete ===');
}

run().catch((err) => {
  console.error('WASM demo error:', err);
});
