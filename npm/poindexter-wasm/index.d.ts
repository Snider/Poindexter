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
// DNS Tools Types
// ============================================================================

/** DNS record types - standard and extended (ClouDNS compatible) */
export type DNSRecordType =
  // Standard record types
  | 'A'
  | 'AAAA'
  | 'MX'
  | 'TXT'
  | 'NS'
  | 'CNAME'
  | 'SOA'
  | 'PTR'
  | 'SRV'
  | 'CAA'
  // Additional record types (ClouDNS and others)
  | 'ALIAS'   // Virtual A record - CNAME-like for apex domain
  | 'RP'      // Responsible Person
  | 'SSHFP'   // SSH Fingerprint
  | 'TLSA'    // DANE TLS Authentication
  | 'DS'      // DNSSEC Delegation Signer
  | 'DNSKEY'  // DNSSEC Key
  | 'NAPTR'   // Naming Authority Pointer
  | 'LOC'     // Geographic Location
  | 'HINFO'   // Host Information
  | 'CERT'    // Certificate record
  | 'SMIMEA'  // S/MIME Certificate Association
  | 'WR'      // Web Redirect (ClouDNS specific)
  | 'SPF';    // Sender Policy Framework (legacy)

/** DNS record type metadata */
export interface DNSRecordTypeInfo {
  type: DNSRecordType;
  name: string;
  description: string;
  rfc?: string;
  common: boolean;
}

/** CAA record */
export interface CAARecord {
  flag: number;
  tag: string;  // "issue", "issuewild", "iodef"
  value: string;
}

/** SSHFP record */
export interface SSHFPRecord {
  algorithm: number;   // 1=RSA, 2=DSA, 3=ECDSA, 4=Ed25519
  fpType: number;      // 1=SHA-1, 2=SHA-256
  fingerprint: string;
}

/** TLSA (DANE) record */
export interface TLSARecord {
  usage: number;        // 0-3: CA constraint, Service cert, Trust anchor, Domain-issued
  selector: number;     // 0=Full cert, 1=SubjectPublicKeyInfo
  matchingType: number; // 0=Exact, 1=SHA-256, 2=SHA-512
  certData: string;
}

/** DS (DNSSEC Delegation Signer) record */
export interface DSRecord {
  keyTag: number;
  algorithm: number;
  digestType: number;
  digest: string;
}

/** DNSKEY record */
export interface DNSKEYRecord {
  flags: number;
  protocol: number;
  algorithm: number;
  publicKey: string;
}

/** NAPTR record */
export interface NAPTRRecord {
  order: number;
  preference: number;
  flags: string;
  service: string;
  regexp: string;
  replacement: string;
}

/** RP (Responsible Person) record */
export interface RPRecord {
  mailbox: string;  // Email as DNS name (user.domain.com)
  txtDom: string;   // Domain with TXT record containing more info
}

/** LOC (Location) record */
export interface LOCRecord {
  latitude: number;
  longitude: number;
  altitude: number;
  size: number;
  hPrecision: number;
  vPrecision: number;
}

/** ALIAS record (provider-specific) */
export interface ALIASRecord {
  target: string;
}

/** Web Redirect record (ClouDNS specific) */
export interface WebRedirectRecord {
  url: string;
  redirectType: number;  // 301, 302, etc.
  frame: boolean;        // Frame redirect vs HTTP redirect
}

/** External tool links for domain/IP/email analysis */
export interface ExternalToolLinks {
  target: string;
  type: 'domain' | 'ip' | 'email';

  // MXToolbox links
  mxtoolboxDns?: string;
  mxtoolboxMx?: string;
  mxtoolboxBlacklist?: string;
  mxtoolboxSmtp?: string;
  mxtoolboxSpf?: string;
  mxtoolboxDmarc?: string;
  mxtoolboxDkim?: string;
  mxtoolboxHttp?: string;
  mxtoolboxHttps?: string;
  mxtoolboxPing?: string;
  mxtoolboxTrace?: string;
  mxtoolboxWhois?: string;
  mxtoolboxAsn?: string;

  // DNSChecker links
  dnscheckerDns?: string;
  dnscheckerPropagation?: string;

  // Other tools
  whois?: string;
  viewdns?: string;
  intodns?: string;
  dnsviz?: string;
  securitytrails?: string;
  shodan?: string;
  censys?: string;
  builtwith?: string;
  ssllabs?: string;
  hstsPreload?: string;
  hardenize?: string;

  // IP-specific tools
  ipinfo?: string;
  abuseipdb?: string;
  virustotal?: string;
  threatcrowd?: string;

  // Email-specific tools
  mailtester?: string;
  learndmarc?: string;
}

/** RDAP server registry */
export interface RDAPServers {
  tlds: Record<string, string>;
  rirs: Record<string, string>;
  universal: string;
}

/** RDAP response event */
export interface RDAPEvent {
  eventAction: string;
  eventDate: string;
  eventActor?: string;
}

/** RDAP entity (registrar, registrant, etc.) */
export interface RDAPEntity {
  handle?: string;
  roles?: string[];
  vcardArray?: any[];
  entities?: RDAPEntity[];
  events?: RDAPEvent[];
}

/** RDAP nameserver */
export interface RDAPNameserver {
  ldhName: string;
  ipAddresses?: {
    v4?: string[];
    v6?: string[];
  };
}

/** RDAP link */
export interface RDAPLink {
  value?: string;
  rel?: string;
  href?: string;
  type?: string;
}

/** RDAP remark/notice */
export interface RDAPRemark {
  title?: string;
  description?: string[];
  links?: RDAPLink[];
}

/** RDAP response (for domain, IP, or ASN lookups) */
export interface RDAPResponse {
  // Common fields
  handle?: string;
  ldhName?: string;
  unicodeName?: string;
  status?: string[];
  events?: RDAPEvent[];
  entities?: RDAPEntity[];
  nameservers?: RDAPNameserver[];
  links?: RDAPLink[];
  remarks?: RDAPRemark[];
  notices?: RDAPRemark[];

  // Network-specific (for IP lookups)
  startAddress?: string;
  endAddress?: string;
  ipVersion?: string;
  name?: string;
  type?: string;
  country?: string;
  parentHandle?: string;

  // Error fields
  errorCode?: number;
  title?: string;
  description?: string[];

  // Metadata
  rawJson?: string;
  lookupTimeMs: number;
  timestamp: string;
  error?: string;
}

/** Parsed domain info from RDAP */
export interface ParsedDomainInfo {
  domain: string;
  registrar?: string;
  registrationDate?: string;
  expirationDate?: string;
  updatedDate?: string;
  status?: string[];
  nameservers?: string[];
  dnssec: boolean;
}

/** DNS lookup result */
export interface DNSLookupResult {
  domain: string;
  queryType: string;
  records: DNSRecord[];
  mxRecords?: MXRecord[];
  srvRecords?: SRVRecord[];
  soaRecord?: SOARecord;
  lookupTimeMs: number;
  error?: string;
  timestamp: string;
}

/** DNS record */
export interface DNSRecord {
  type: DNSRecordType;
  name: string;
  value: string;
  ttl?: number;
}

/** MX record */
export interface MXRecord {
  host: string;
  priority: number;
}

/** SRV record */
export interface SRVRecord {
  target: string;
  port: number;
  priority: number;
  weight: number;
}

/** SOA record */
export interface SOARecord {
  primaryNs: string;
  adminEmail: string;
  serial: number;
  refresh: number;
  retry: number;
  expire: number;
  minTtl: number;
}

/** Complete DNS lookup result */
export interface CompleteDNSLookup {
  domain: string;
  a?: string[];
  aaaa?: string[];
  mx?: MXRecord[];
  ns?: string[];
  txt?: string[];
  cname?: string;
  soa?: SOARecord;
  lookupTimeMs: number;
  errors?: string[];
  timestamp: string;
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

  // DNS tools
  getExternalToolLinks(domain: string): Promise<ExternalToolLinks>;
  getExternalToolLinksIP(ip: string): Promise<ExternalToolLinks>;
  getExternalToolLinksEmail(emailOrDomain: string): Promise<ExternalToolLinks>;
  getRDAPServers(): Promise<RDAPServers>;
  buildRDAPDomainURL(domain: string): Promise<string>;
  buildRDAPIPURL(ip: string): Promise<string>;
  buildRDAPASNURL(asn: string): Promise<string>;
  getDNSRecordTypes(): Promise<DNSRecordType[]>;
  getDNSRecordTypeInfo(): Promise<DNSRecordTypeInfo[]>;
  getCommonDNSRecordTypes(): Promise<DNSRecordType[]>;
}

export function init(options?: InitOptions): Promise<PxAPI>;
