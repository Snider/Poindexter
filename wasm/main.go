//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	pd "github.com/Snider/Poindexter"
)

// Simple registry for KDTree instances created from JS.
// We keep values as string for simplicity across the WASM boundary.
var (
	treeRegistry = map[int]*pd.KDTree[string]{}
	nextTreeID   = 1
)

func export(name string, fn func(this js.Value, args []js.Value) (any, error)) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		res, err := fn(this, args)
		if err != nil {
			return map[string]any{"ok": false, "error": err.Error()}
		}
		return map[string]any{"ok": true, "data": res}
	}))
}

func getInt(v js.Value, idx int) (int, error) {
	if len := v.Length(); len > idx {
		return v.Index(idx).Int(), nil
	}
	return 0, errors.New("missing integer argument")
}

func getFloatSlice(arg js.Value) ([]float64, error) {
	if arg.IsUndefined() || arg.IsNull() {
		return nil, errors.New("coords/query is undefined or null")
	}
	ln := arg.Length()
	res := make([]float64, ln)
	for i := 0; i < ln; i++ {
		res[i] = arg.Index(i).Float()
	}
	return res, nil
}

func version(_ js.Value, _ []js.Value) (any, error) {
	return pd.Version(), nil
}

func hello(_ js.Value, args []js.Value) (any, error) {
	name := ""
	if len(args) > 0 {
		name = args[0].String()
	}
	return pd.Hello(name), nil
}

func newTree(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("newTree(dim) requires dim")
	}
	dim := args[0].Int()
	if dim <= 0 {
		return nil, pd.ErrZeroDim
	}
	t, err := pd.NewKDTreeFromDim[string](dim)
	if err != nil {
		return nil, err
	}
	id := nextTreeID
	nextTreeID++
	treeRegistry[id] = t
	return map[string]any{"treeId": id, "dim": dim}, nil
}

func treeLen(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("len(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.Len(), nil
}

func treeDim(_ js.Value, args []js.Value) (any, error) {
	if len(args) < 1 {
		return nil, errors.New("dim(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.Dim(), nil
}

func insert(_ js.Value, args []js.Value) (any, error) {
	// insert(treeId, {id: string, coords: number[], value?: string})
	if len(args) < 2 {
		return nil, errors.New("insert(treeId, point)")
	}
	id := args[0].Int()
	pt := args[1]
	pid := pt.Get("id").String()
	coords, err := getFloatSlice(pt.Get("coords"))
	if err != nil {
		return nil, err
	}
	val := pt.Get("value").String()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	okIns := t.Insert(pd.KDPoint[string]{ID: pid, Coords: coords, Value: val})
	return okIns, nil
}

func deleteByID(_ js.Value, args []js.Value) (any, error) {
	// deleteByID(treeId, id)
	if len(args) < 2 {
		return nil, errors.New("deleteByID(treeId, id)")
	}
	id := args[0].Int()
	pid := args[1].String()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	return t.DeleteByID(pid), nil
}

func nearest(_ js.Value, args []js.Value) (any, error) {
	// nearest(treeId, query:number[]) -> {point, dist, found}
	if len(args) < 2 {
		return nil, errors.New("nearest(treeId, query)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	p, d, found := t.Nearest(query)
	out := map[string]any{
		"point": map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value},
		"dist":  d,
		"found": found,
	}
	return out, nil
}

func kNearest(_ js.Value, args []js.Value) (any, error) {
	// kNearest(treeId, query:number[], k:int) -> {points:[...], dists:[...]}
	if len(args) < 3 {
		return nil, errors.New("kNearest(treeId, query, k)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	k := args[2].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	pts, dists := t.KNearest(query, k)
	jsPts := make([]any, len(pts))
	for i, p := range pts {
		jsPts[i] = map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value}
	}
	return map[string]any{"points": jsPts, "dists": dists}, nil
}

func radius(_ js.Value, args []js.Value) (any, error) {
	// radius(treeId, query:number[], r:number) -> {points:[...], dists:[...]}
	if len(args) < 3 {
		return nil, errors.New("radius(treeId, query, r)")
	}
	id := args[0].Int()
	query, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}
	r := args[2].Float()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	pts, dists := t.Radius(query, r)
	jsPts := make([]any, len(pts))
	for i, p := range pts {
		jsPts[i] = map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value}
	}
	return map[string]any{"points": jsPts, "dists": dists}, nil
}

func exportJSON(_ js.Value, args []js.Value) (any, error) {
	// exportJSON(treeId) -> string (all points)
	if len(args) < 1 {
		return nil, errors.New("exportJSON(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	// Export all points
	points := t.Points()
	jsPts := make([]any, len(points))
	for i, p := range points {
		jsPts[i] = map[string]any{"id": p.ID, "coords": p.Coords, "value": p.Value}
	}
	m := map[string]any{
		"dim":     t.Dim(),
		"len":     t.Len(),
		"backend": string(t.Backend()),
		"points":  jsPts,
	}
	b, _ := json.Marshal(m)
	return string(b), nil
}

func getAnalytics(_ js.Value, args []js.Value) (any, error) {
	// getAnalytics(treeId) -> analytics snapshot
	if len(args) < 1 {
		return nil, errors.New("getAnalytics(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	snap := t.GetAnalyticsSnapshot()
	return map[string]any{
		"queryCount":          snap.QueryCount,
		"insertCount":         snap.InsertCount,
		"deleteCount":         snap.DeleteCount,
		"avgQueryTimeNs":      snap.AvgQueryTimeNs,
		"minQueryTimeNs":      snap.MinQueryTimeNs,
		"maxQueryTimeNs":      snap.MaxQueryTimeNs,
		"lastQueryTimeNs":     snap.LastQueryTimeNs,
		"lastQueryAt":         snap.LastQueryAt.UnixMilli(),
		"createdAt":           snap.CreatedAt.UnixMilli(),
		"backendRebuildCount": snap.BackendRebuildCnt,
		"lastRebuiltAt":       snap.LastRebuiltAt.UnixMilli(),
	}, nil
}

func getPeerStats(_ js.Value, args []js.Value) (any, error) {
	// getPeerStats(treeId) -> array of peer stats
	if len(args) < 1 {
		return nil, errors.New("getPeerStats(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	stats := t.GetPeerStats()
	jsStats := make([]any, len(stats))
	for i, s := range stats {
		jsStats[i] = map[string]any{
			"peerId":         s.PeerID,
			"selectionCount": s.SelectionCount,
			"avgDistance":    s.AvgDistance,
			"lastSelectedAt": s.LastSelectedAt.UnixMilli(),
		}
	}
	return jsStats, nil
}

func getTopPeers(_ js.Value, args []js.Value) (any, error) {
	// getTopPeers(treeId, n) -> array of top n peer stats
	if len(args) < 2 {
		return nil, errors.New("getTopPeers(treeId, n)")
	}
	id := args[0].Int()
	n := args[1].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	stats := t.GetTopPeers(n)
	jsStats := make([]any, len(stats))
	for i, s := range stats {
		jsStats[i] = map[string]any{
			"peerId":         s.PeerID,
			"selectionCount": s.SelectionCount,
			"avgDistance":    s.AvgDistance,
			"lastSelectedAt": s.LastSelectedAt.UnixMilli(),
		}
	}
	return jsStats, nil
}

func getAxisDistributions(_ js.Value, args []js.Value) (any, error) {
	// getAxisDistributions(treeId, axisNames?: string[]) -> array of axis distribution stats
	if len(args) < 1 {
		return nil, errors.New("getAxisDistributions(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}

	var axisNames []string
	if len(args) > 1 && !args[1].IsUndefined() && !args[1].IsNull() {
		ln := args[1].Length()
		axisNames = make([]string, ln)
		for i := 0; i < ln; i++ {
			axisNames[i] = args[1].Index(i).String()
		}
	}

	dists := t.ComputeDistanceDistribution(axisNames)
	jsDists := make([]any, len(dists))
	for i, d := range dists {
		jsDists[i] = map[string]any{
			"axis": d.Axis,
			"name": d.Name,
			"stats": map[string]any{
				"count":    d.Stats.Count,
				"min":      d.Stats.Min,
				"max":      d.Stats.Max,
				"mean":     d.Stats.Mean,
				"median":   d.Stats.Median,
				"stdDev":   d.Stats.StdDev,
				"p25":      d.Stats.P25,
				"p75":      d.Stats.P75,
				"p90":      d.Stats.P90,
				"p99":      d.Stats.P99,
				"variance": d.Stats.Variance,
				"skewness": d.Stats.Skewness,
			},
		}
	}
	return jsDists, nil
}

func resetAnalytics(_ js.Value, args []js.Value) (any, error) {
	// resetAnalytics(treeId) -> resets all analytics
	if len(args) < 1 {
		return nil, errors.New("resetAnalytics(treeId)")
	}
	id := args[0].Int()
	t, ok := treeRegistry[id]
	if !ok {
		return nil, fmt.Errorf("unknown treeId %d", id)
	}
	t.ResetAnalytics()
	return true, nil
}

func computeDistributionStats(_ js.Value, args []js.Value) (any, error) {
	// computeDistributionStats(distances: number[]) -> distribution stats
	if len(args) < 1 {
		return nil, errors.New("computeDistributionStats(distances)")
	}
	distances, err := getFloatSlice(args[0])
	if err != nil {
		return nil, err
	}
	stats := pd.ComputeDistributionStats(distances)
	return map[string]any{
		"count":      stats.Count,
		"min":        stats.Min,
		"max":        stats.Max,
		"mean":       stats.Mean,
		"median":     stats.Median,
		"stdDev":     stats.StdDev,
		"p25":        stats.P25,
		"p75":        stats.P75,
		"p90":        stats.P90,
		"p99":        stats.P99,
		"variance":   stats.Variance,
		"skewness":   stats.Skewness,
		"sampleSize": stats.SampleSize,
		"computedAt": stats.ComputedAt.UnixMilli(),
	}, nil
}

func computePeerQualityScore(_ js.Value, args []js.Value) (any, error) {
	// computePeerQualityScore(metrics: NATRoutingMetrics, weights?: QualityWeights) -> score
	if len(args) < 1 {
		return nil, errors.New("computePeerQualityScore(metrics)")
	}
	m := args[0]
	metrics := pd.NATRoutingMetrics{
		ConnectivityScore: m.Get("connectivityScore").Float(),
		SymmetryScore:     m.Get("symmetryScore").Float(),
		RelayProbability:  m.Get("relayProbability").Float(),
		DirectSuccessRate: m.Get("directSuccessRate").Float(),
		AvgRTTMs:          m.Get("avgRttMs").Float(),
		JitterMs:          m.Get("jitterMs").Float(),
		PacketLossRate:    m.Get("packetLossRate").Float(),
		BandwidthMbps:     m.Get("bandwidthMbps").Float(),
		NATType:           m.Get("natType").String(),
	}

	var weights *pd.QualityWeights
	if len(args) > 1 && !args[1].IsUndefined() && !args[1].IsNull() {
		w := args[1]
		weights = &pd.QualityWeights{
			Latency:       w.Get("latency").Float(),
			Jitter:        w.Get("jitter").Float(),
			PacketLoss:    w.Get("packetLoss").Float(),
			Bandwidth:     w.Get("bandwidth").Float(),
			Connectivity:  w.Get("connectivity").Float(),
			Symmetry:      w.Get("symmetry").Float(),
			DirectSuccess: w.Get("directSuccess").Float(),
			RelayPenalty:  w.Get("relayPenalty").Float(),
			NATType:       w.Get("natType").Float(),
		}
	}

	score := pd.PeerQualityScore(metrics, weights)
	return score, nil
}

func computeTrustScore(_ js.Value, args []js.Value) (any, error) {
	// computeTrustScore(metrics: TrustMetrics) -> score
	if len(args) < 1 {
		return nil, errors.New("computeTrustScore(metrics)")
	}
	m := args[0]
	metrics := pd.TrustMetrics{
		ReputationScore:        m.Get("reputationScore").Float(),
		SuccessfulTransactions: int64(m.Get("successfulTransactions").Int()),
		FailedTransactions:     int64(m.Get("failedTransactions").Int()),
		AgeSeconds:             int64(m.Get("ageSeconds").Int()),
		VouchCount:             m.Get("vouchCount").Int(),
		FlagCount:              m.Get("flagCount").Int(),
		ProofOfWork:            m.Get("proofOfWork").Float(),
	}

	score := pd.ComputeTrustScore(metrics)
	return score, nil
}

func getDefaultQualityWeights(_ js.Value, _ []js.Value) (any, error) {
	w := pd.DefaultQualityWeights()
	return map[string]any{
		"latency":       w.Latency,
		"jitter":        w.Jitter,
		"packetLoss":    w.PacketLoss,
		"bandwidth":     w.Bandwidth,
		"connectivity":  w.Connectivity,
		"symmetry":      w.Symmetry,
		"directSuccess": w.DirectSuccess,
		"relayPenalty":  w.RelayPenalty,
		"natType":       w.NATType,
	}, nil
}

func getDefaultPeerFeatureRanges(_ js.Value, _ []js.Value) (any, error) {
	ranges := pd.DefaultPeerFeatureRanges()
	jsRanges := make([]any, len(ranges.Ranges))
	for i, r := range ranges.Ranges {
		jsRanges[i] = map[string]any{
			"min": r.Min,
			"max": r.Max,
		}
	}
	return map[string]any{
		"ranges": jsRanges,
		"labels": pd.StandardFeatureLabels(),
	}, nil
}

func normalizePeerFeatures(_ js.Value, args []js.Value) (any, error) {
	// normalizePeerFeatures(features: number[], ranges?: FeatureRanges) -> number[]
	if len(args) < 1 {
		return nil, errors.New("normalizePeerFeatures(features)")
	}
	features, err := getFloatSlice(args[0])
	if err != nil {
		return nil, err
	}

	ranges := pd.DefaultPeerFeatureRanges()
	if len(args) > 1 && !args[1].IsUndefined() && !args[1].IsNull() {
		rangesArg := args[1].Get("ranges")
		if !rangesArg.IsUndefined() && !rangesArg.IsNull() {
			ln := rangesArg.Length()
			ranges.Ranges = make([]pd.AxisStats, ln)
			for i := 0; i < ln; i++ {
				r := rangesArg.Index(i)
				ranges.Ranges[i] = pd.AxisStats{
					Min: r.Get("min").Float(),
					Max: r.Get("max").Float(),
				}
			}
		}
	}

	normalized := pd.NormalizePeerFeatures(features, ranges)
	return normalized, nil
}

func weightedPeerFeatures(_ js.Value, args []js.Value) (any, error) {
	// weightedPeerFeatures(normalized: number[], weights: number[]) -> number[]
	if len(args) < 2 {
		return nil, errors.New("weightedPeerFeatures(normalized, weights)")
	}
	normalized, err := getFloatSlice(args[0])
	if err != nil {
		return nil, err
	}
	weights, err := getFloatSlice(args[1])
	if err != nil {
		return nil, err
	}

	weighted := pd.WeightedPeerFeatures(normalized, weights)
	return weighted, nil
}

// ============================================================================
// DNS Tools Functions
// ============================================================================

func getExternalToolLinks(_ js.Value, args []js.Value) (any, error) {
	// getExternalToolLinks(domain: string) -> ExternalToolLinks
	if len(args) < 1 {
		return nil, errors.New("getExternalToolLinks(domain)")
	}
	domain := args[0].String()
	links := pd.GetExternalToolLinks(domain)
	return externalToolLinksToJS(links), nil
}

func getExternalToolLinksIP(_ js.Value, args []js.Value) (any, error) {
	// getExternalToolLinksIP(ip: string) -> ExternalToolLinks
	if len(args) < 1 {
		return nil, errors.New("getExternalToolLinksIP(ip)")
	}
	ip := args[0].String()
	links := pd.GetExternalToolLinksIP(ip)
	return externalToolLinksToJS(links), nil
}

func getExternalToolLinksEmail(_ js.Value, args []js.Value) (any, error) {
	// getExternalToolLinksEmail(emailOrDomain: string) -> ExternalToolLinks
	if len(args) < 1 {
		return nil, errors.New("getExternalToolLinksEmail(emailOrDomain)")
	}
	emailOrDomain := args[0].String()
	links := pd.GetExternalToolLinksEmail(emailOrDomain)
	return externalToolLinksToJS(links), nil
}

func externalToolLinksToJS(links pd.ExternalToolLinks) map[string]any {
	return map[string]any{
		"target": links.Target,
		"type":   links.Type,
		// MXToolbox
		"mxtoolboxDns":       links.MXToolboxDNS,
		"mxtoolboxMx":        links.MXToolboxMX,
		"mxtoolboxBlacklist": links.MXToolboxBlacklist,
		"mxtoolboxSmtp":      links.MXToolboxSMTP,
		"mxtoolboxSpf":       links.MXToolboxSPF,
		"mxtoolboxDmarc":     links.MXToolboxDMARC,
		"mxtoolboxDkim":      links.MXToolboxDKIM,
		"mxtoolboxHttp":      links.MXToolboxHTTP,
		"mxtoolboxHttps":     links.MXToolboxHTTPS,
		"mxtoolboxPing":      links.MXToolboxPing,
		"mxtoolboxTrace":     links.MXToolboxTrace,
		"mxtoolboxWhois":     links.MXToolboxWhois,
		"mxtoolboxAsn":       links.MXToolboxASN,
		// DNSChecker
		"dnscheckerDns":         links.DNSCheckerDNS,
		"dnscheckerPropagation": links.DNSCheckerPropagation,
		// Other tools
		"whois":          links.WhoIs,
		"viewdns":        links.ViewDNS,
		"intodns":        links.IntoDNS,
		"dnsviz":         links.DNSViz,
		"securitytrails": links.SecurityTrails,
		"shodan":         links.Shodan,
		"censys":         links.Censys,
		"builtwith":      links.BuiltWith,
		"ssllabs":        links.SSLLabs,
		"hstsPreload":    links.HSTSPreload,
		"hardenize":      links.Hardenize,
		// IP-specific
		"ipinfo":      links.IPInfo,
		"abuseipdb":   links.AbuseIPDB,
		"virustotal":  links.VirusTotal,
		"threatcrowd": links.ThreatCrowd,
		// Email-specific
		"mailtester": links.MailTester,
		"learndmarc": links.LearnDMARC,
	}
}

func getRDAPServers(_ js.Value, _ []js.Value) (any, error) {
	// Returns a list of known RDAP servers for reference
	servers := map[string]any{
		"tlds": map[string]string{
			"com":  "https://rdap.verisign.com/com/v1/",
			"net":  "https://rdap.verisign.com/net/v1/",
			"org":  "https://rdap.publicinterestregistry.org/rdap/",
			"info": "https://rdap.afilias.net/rdap/info/",
			"io":   "https://rdap.nic.io/",
			"co":   "https://rdap.nic.co/",
			"dev":  "https://rdap.nic.google/",
			"app":  "https://rdap.nic.google/",
		},
		"rirs": map[string]string{
			"arin":    "https://rdap.arin.net/registry/",
			"ripe":    "https://rdap.db.ripe.net/",
			"apnic":   "https://rdap.apnic.net/",
			"afrinic": "https://rdap.afrinic.net/rdap/",
			"lacnic":  "https://rdap.lacnic.net/rdap/",
		},
		"universal": "https://rdap.org/",
	}
	return servers, nil
}

func buildRDAPDomainURL(_ js.Value, args []js.Value) (any, error) {
	// buildRDAPDomainURL(domain: string) -> string
	if len(args) < 1 {
		return nil, errors.New("buildRDAPDomainURL(domain)")
	}
	domain := args[0].String()
	// Use universal RDAP redirector
	return fmt.Sprintf("https://rdap.org/domain/%s", domain), nil
}

func buildRDAPIPURL(_ js.Value, args []js.Value) (any, error) {
	// buildRDAPIPURL(ip: string) -> string
	if len(args) < 1 {
		return nil, errors.New("buildRDAPIPURL(ip)")
	}
	ip := args[0].String()
	return fmt.Sprintf("https://rdap.org/ip/%s", ip), nil
}

func buildRDAPASNURL(_ js.Value, args []js.Value) (any, error) {
	// buildRDAPASNURL(asn: string) -> string
	if len(args) < 1 {
		return nil, errors.New("buildRDAPASNURL(asn)")
	}
	asn := args[0].String()
	// Normalize ASN
	asnNum := asn
	if len(asn) > 2 && (asn[:2] == "AS" || asn[:2] == "as") {
		asnNum = asn[2:]
	}
	return fmt.Sprintf("https://rdap.org/autnum/%s", asnNum), nil
}

func getDNSRecordTypes(_ js.Value, _ []js.Value) (any, error) {
	// Returns all available DNS record types
	types := pd.GetAllDNSRecordTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result, nil
}

func getDNSRecordTypeInfo(_ js.Value, _ []js.Value) (any, error) {
	// Returns detailed info about all DNS record types
	info := pd.GetDNSRecordTypeInfo()
	result := make([]any, len(info))
	for i, r := range info {
		result[i] = map[string]any{
			"type":        string(r.Type),
			"name":        r.Name,
			"description": r.Description,
			"rfc":         r.RFC,
			"common":      r.Common,
		}
	}
	return result, nil
}

func getCommonDNSRecordTypes(_ js.Value, _ []js.Value) (any, error) {
	// Returns only commonly used DNS record types
	types := pd.GetCommonDNSRecordTypes()
	result := make([]string, len(types))
	for i, t := range types {
		result[i] = string(t)
	}
	return result, nil
}

func main() {
	// Export core API
	export("pxVersion", version)
	export("pxHello", hello)
	export("pxNewTree", newTree)
	export("pxTreeLen", treeLen)
	export("pxTreeDim", treeDim)
	export("pxInsert", insert)
	export("pxDeleteByID", deleteByID)
	export("pxNearest", nearest)
	export("pxKNearest", kNearest)
	export("pxRadius", radius)
	export("pxExportJSON", exportJSON)

	// Export analytics API
	export("pxGetAnalytics", getAnalytics)
	export("pxGetPeerStats", getPeerStats)
	export("pxGetTopPeers", getTopPeers)
	export("pxGetAxisDistributions", getAxisDistributions)
	export("pxResetAnalytics", resetAnalytics)
	export("pxComputeDistributionStats", computeDistributionStats)

	// Export NAT routing / peer quality API
	export("pxComputePeerQualityScore", computePeerQualityScore)
	export("pxComputeTrustScore", computeTrustScore)
	export("pxGetDefaultQualityWeights", getDefaultQualityWeights)
	export("pxGetDefaultPeerFeatureRanges", getDefaultPeerFeatureRanges)
	export("pxNormalizePeerFeatures", normalizePeerFeatures)
	export("pxWeightedPeerFeatures", weightedPeerFeatures)

	// Export DNS tools API
	export("pxGetExternalToolLinks", getExternalToolLinks)
	export("pxGetExternalToolLinksIP", getExternalToolLinksIP)
	export("pxGetExternalToolLinksEmail", getExternalToolLinksEmail)
	export("pxGetRDAPServers", getRDAPServers)
	export("pxBuildRDAPDomainURL", buildRDAPDomainURL)
	export("pxBuildRDAPIPURL", buildRDAPIPURL)
	export("pxBuildRDAPASNURL", buildRDAPASNURL)
	export("pxGetDNSRecordTypes", getDNSRecordTypes)
	export("pxGetDNSRecordTypeInfo", getDNSRecordTypeInfo)
	export("pxGetCommonDNSRecordTypes", getCommonDNSRecordTypes)

	// Keep running
	select {}
}
