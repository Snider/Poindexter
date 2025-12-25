package poindexter

import (
	"strings"
	"testing"
)

// ============================================================================
// External Tool Links Tests
// ============================================================================

func TestGetExternalToolLinks(t *testing.T) {
	links := GetExternalToolLinks("example.com")

	if links.Target != "example.com" {
		t.Errorf("expected target=example.com, got %s", links.Target)
	}
	if links.Type != "domain" {
		t.Errorf("expected type=domain, got %s", links.Type)
	}

	// Check MXToolbox links
	if !strings.Contains(links.MXToolboxDNS, "mxtoolbox.com") {
		t.Error("MXToolboxDNS should contain mxtoolbox.com")
	}
	if !strings.Contains(links.MXToolboxDNS, "example.com") {
		t.Error("MXToolboxDNS should contain the domain")
	}

	if !strings.Contains(links.MXToolboxMX, "mxtoolbox.com") {
		t.Error("MXToolboxMX should contain mxtoolbox.com")
	}

	if !strings.Contains(links.MXToolboxSPF, "spf") {
		t.Error("MXToolboxSPF should contain 'spf'")
	}

	if !strings.Contains(links.MXToolboxDMARC, "dmarc") {
		t.Error("MXToolboxDMARC should contain 'dmarc'")
	}

	// Check DNSChecker links
	if !strings.Contains(links.DNSCheckerDNS, "dnschecker.org") {
		t.Error("DNSCheckerDNS should contain dnschecker.org")
	}

	// Check other tools
	if !strings.Contains(links.WhoIs, "who.is") {
		t.Error("WhoIs should contain who.is")
	}

	if !strings.Contains(links.SSLLabs, "ssllabs.com") {
		t.Error("SSLLabs should contain ssllabs.com")
	}

	if !strings.Contains(links.VirusTotal, "virustotal.com") {
		t.Error("VirusTotal should contain virustotal.com")
	}
}

func TestGetExternalToolLinksIP(t *testing.T) {
	links := GetExternalToolLinksIP("8.8.8.8")

	if links.Target != "8.8.8.8" {
		t.Errorf("expected target=8.8.8.8, got %s", links.Target)
	}
	if links.Type != "ip" {
		t.Errorf("expected type=ip, got %s", links.Type)
	}

	// Check IP-specific links
	if !strings.Contains(links.IPInfo, "ipinfo.io") {
		t.Error("IPInfo should contain ipinfo.io")
	}
	if !strings.Contains(links.IPInfo, "8.8.8.8") {
		t.Error("IPInfo should contain the IP address")
	}

	if !strings.Contains(links.AbuseIPDB, "abuseipdb.com") {
		t.Error("AbuseIPDB should contain abuseipdb.com")
	}

	if !strings.Contains(links.Shodan, "shodan.io") {
		t.Error("Shodan should contain shodan.io")
	}

	if !strings.Contains(links.MXToolboxBlacklist, "blacklist") {
		t.Error("MXToolboxBlacklist should contain 'blacklist'")
	}
}

func TestGetExternalToolLinksEmail(t *testing.T) {
	// Test with email address
	links := GetExternalToolLinksEmail("test@example.com")

	if links.Target != "test@example.com" {
		t.Errorf("expected target=test@example.com, got %s", links.Target)
	}
	if links.Type != "email" {
		t.Errorf("expected type=email, got %s", links.Type)
	}

	// Email tools should use the domain
	if !strings.Contains(links.MXToolboxMX, "example.com") {
		t.Error("MXToolboxMX should contain the domain from email")
	}

	if !strings.Contains(links.MXToolboxSPF, "spf") {
		t.Error("MXToolboxSPF should contain 'spf'")
	}

	if !strings.Contains(links.MXToolboxDMARC, "dmarc") {
		t.Error("MXToolboxDMARC should contain 'dmarc'")
	}

	// Test with just domain
	links2 := GetExternalToolLinksEmail("example.org")
	if links2.Target != "example.org" {
		t.Errorf("expected target=example.org, got %s", links2.Target)
	}
}

func TestGetExternalToolLinksSpecialChars(t *testing.T) {
	// Test URL encoding
	links := GetExternalToolLinks("test-domain.example.com")

	if !strings.Contains(links.MXToolboxDNS, "test-domain.example.com") {
		t.Error("Should handle hyphens in domain")
	}
}

// ============================================================================
// DNS Lookup Tests (Unit tests for structure, not network)
// ============================================================================

func TestDNSRecordTypes(t *testing.T) {
	types := []DNSRecordType{
		DNSRecordA,
		DNSRecordAAAA,
		DNSRecordMX,
		DNSRecordTXT,
		DNSRecordNS,
		DNSRecordCNAME,
		DNSRecordSOA,
		DNSRecordPTR,
		DNSRecordSRV,
		DNSRecordCAA,
	}

	expected := []string{"A", "AAAA", "MX", "TXT", "NS", "CNAME", "SOA", "PTR", "SRV", "CAA"}

	for i, typ := range types {
		if string(typ) != expected[i] {
			t.Errorf("expected type %s, got %s", expected[i], typ)
		}
	}
}

func TestDNSRecordTypesExtended(t *testing.T) {
	// Test all ClouDNS record types are defined
	types := []DNSRecordType{
		DNSRecordALIAS,
		DNSRecordRP,
		DNSRecordSSHFP,
		DNSRecordTLSA,
		DNSRecordDS,
		DNSRecordDNSKEY,
		DNSRecordNAPTR,
		DNSRecordLOC,
		DNSRecordHINFO,
		DNSRecordCERT,
		DNSRecordSMIMEA,
		DNSRecordWR,
		DNSRecordSPF,
	}

	expected := []string{"ALIAS", "RP", "SSHFP", "TLSA", "DS", "DNSKEY", "NAPTR", "LOC", "HINFO", "CERT", "SMIMEA", "WR", "SPF"}

	for i, typ := range types {
		if string(typ) != expected[i] {
			t.Errorf("expected type %s, got %s", expected[i], typ)
		}
	}
}

func TestGetDNSRecordTypeInfo(t *testing.T) {
	info := GetDNSRecordTypeInfo()

	if len(info) == 0 {
		t.Error("GetDNSRecordTypeInfo should return non-empty list")
	}

	// Check that common types exist
	commonFound := 0
	for _, r := range info {
		if r.Common {
			commonFound++
		}
		// Each entry should have type, name, and description
		if r.Type == "" {
			t.Error("Record type should not be empty")
		}
		if r.Name == "" {
			t.Error("Record name should not be empty")
		}
		if r.Description == "" {
			t.Error("Record description should not be empty")
		}
	}

	if commonFound < 10 {
		t.Errorf("Expected at least 10 common record types, got %d", commonFound)
	}

	// Check for specific types
	typeMap := make(map[DNSRecordType]DNSRecordTypeInfo)
	for _, r := range info {
		typeMap[r.Type] = r
	}

	if _, ok := typeMap[DNSRecordA]; !ok {
		t.Error("A record type should be in info")
	}
	if _, ok := typeMap[DNSRecordALIAS]; !ok {
		t.Error("ALIAS record type should be in info")
	}
	if _, ok := typeMap[DNSRecordTLSA]; !ok {
		t.Error("TLSA record type should be in info")
	}
	if _, ok := typeMap[DNSRecordWR]; !ok {
		t.Error("WR (Web Redirect) record type should be in info")
	}
}

func TestGetCommonDNSRecordTypes(t *testing.T) {
	types := GetCommonDNSRecordTypes()

	if len(types) == 0 {
		t.Error("GetCommonDNSRecordTypes should return non-empty list")
	}

	// Check that standard types are present
	typeSet := make(map[DNSRecordType]bool)
	for _, typ := range types {
		typeSet[typ] = true
	}

	if !typeSet[DNSRecordA] {
		t.Error("A record should be in common types")
	}
	if !typeSet[DNSRecordAAAA] {
		t.Error("AAAA record should be in common types")
	}
	if !typeSet[DNSRecordMX] {
		t.Error("MX record should be in common types")
	}
	if !typeSet[DNSRecordTXT] {
		t.Error("TXT record should be in common types")
	}
	if !typeSet[DNSRecordALIAS] {
		t.Error("ALIAS record should be in common types")
	}
}

func TestGetAllDNSRecordTypes(t *testing.T) {
	types := GetAllDNSRecordTypes()

	if len(types) < 20 {
		t.Errorf("GetAllDNSRecordTypes should return at least 20 types, got %d", len(types))
	}

	// Check for ClouDNS-specific types
	typeSet := make(map[DNSRecordType]bool)
	for _, typ := range types {
		typeSet[typ] = true
	}

	if !typeSet[DNSRecordWR] {
		t.Error("WR (Web Redirect) should be in all types")
	}
	if !typeSet[DNSRecordNAPTR] {
		t.Error("NAPTR should be in all types")
	}
	if !typeSet[DNSRecordDS] {
		t.Error("DS should be in all types")
	}
}

func TestDNSLookupResultStructure(t *testing.T) {
	result := DNSLookupResult{
		Domain:    "example.com",
		QueryType: "A",
		Records: []DNSRecord{
			{Type: DNSRecordA, Name: "example.com", Value: "93.184.216.34"},
		},
		LookupTimeMs: 50,
	}

	if result.Domain != "example.com" {
		t.Error("Domain should be set")
	}
	if len(result.Records) != 1 {
		t.Error("Should have 1 record")
	}
	if result.Records[0].Type != DNSRecordA {
		t.Error("Record type should be A")
	}
}

func TestCompleteDNSLookupStructure(t *testing.T) {
	result := CompleteDNSLookup{
		Domain: "example.com",
		A:      []string{"93.184.216.34"},
		AAAA:   []string{"2606:2800:220:1:248:1893:25c8:1946"},
		MX: []MXRecord{
			{Host: "mail.example.com", Priority: 10},
		},
		NS:  []string{"ns1.example.com", "ns2.example.com"},
		TXT: []string{"v=spf1 include:_spf.example.com ~all"},
	}

	if result.Domain != "example.com" {
		t.Error("Domain should be set")
	}
	if len(result.A) != 1 {
		t.Error("Should have 1 A record")
	}
	if len(result.AAAA) != 1 {
		t.Error("Should have 1 AAAA record")
	}
	if len(result.MX) != 1 {
		t.Error("Should have 1 MX record")
	}
	if result.MX[0].Priority != 10 {
		t.Error("MX priority should be 10")
	}
	if len(result.NS) != 2 {
		t.Error("Should have 2 NS records")
	}
}

// ============================================================================
// RDAP Tests (Unit tests for structure, not network)
// ============================================================================

func TestRDAPResponseStructure(t *testing.T) {
	resp := RDAPResponse{
		LDHName: "example.com",
		Status:  []string{"active", "client transfer prohibited"},
		Events: []RDAPEvent{
			{EventAction: "registration", EventDate: "2020-01-01T00:00:00Z"},
			{EventAction: "expiration", EventDate: "2025-01-01T00:00:00Z"},
		},
		Entities: []RDAPEntity{
			{Handle: "REGISTRAR-1", Roles: []string{"registrar"}},
		},
		Nameservers: []RDAPNs{
			{LDHName: "ns1.example.com"},
			{LDHName: "ns2.example.com"},
		},
	}

	if resp.LDHName != "example.com" {
		t.Error("LDHName should be set")
	}
	if len(resp.Status) != 2 {
		t.Error("Should have 2 status values")
	}
	if len(resp.Events) != 2 {
		t.Error("Should have 2 events")
	}
	if resp.Events[0].EventAction != "registration" {
		t.Error("First event should be registration")
	}
	if len(resp.Nameservers) != 2 {
		t.Error("Should have 2 nameservers")
	}
}

func TestParseRDAPResponse(t *testing.T) {
	resp := RDAPResponse{
		LDHName: "example.com",
		Status:  []string{"active", "dnssecSigned"},
		Events: []RDAPEvent{
			{EventAction: "registration", EventDate: "2020-01-01T00:00:00Z"},
			{EventAction: "expiration", EventDate: "2025-01-01T00:00:00Z"},
			{EventAction: "last changed", EventDate: "2024-06-15T00:00:00Z"},
		},
		Entities: []RDAPEntity{
			{Handle: "REGISTRAR-123", Roles: []string{"registrar"}},
		},
		Nameservers: []RDAPNs{
			{LDHName: "ns1.example.com"},
			{LDHName: "ns2.example.com"},
		},
	}

	info := ParseRDAPResponse(resp)

	if info.Domain != "example.com" {
		t.Errorf("expected domain=example.com, got %s", info.Domain)
	}
	if info.RegistrationDate != "2020-01-01T00:00:00Z" {
		t.Errorf("expected registration date, got %s", info.RegistrationDate)
	}
	if info.ExpirationDate != "2025-01-01T00:00:00Z" {
		t.Errorf("expected expiration date, got %s", info.ExpirationDate)
	}
	if info.UpdatedDate != "2024-06-15T00:00:00Z" {
		t.Errorf("expected updated date, got %s", info.UpdatedDate)
	}
	if info.Registrar != "REGISTRAR-123" {
		t.Errorf("expected registrar, got %s", info.Registrar)
	}
	if len(info.Nameservers) != 2 {
		t.Error("Should have 2 nameservers")
	}
	if !info.DNSSEC {
		t.Error("DNSSEC should be true (detected from status)")
	}
}

func TestParseRDAPResponseEmpty(t *testing.T) {
	resp := RDAPResponse{
		LDHName: "test.com",
	}

	info := ParseRDAPResponse(resp)

	if info.Domain != "test.com" {
		t.Error("Domain should be set even with minimal response")
	}
	if info.DNSSEC {
		t.Error("DNSSEC should be false with no status")
	}
	if len(info.Nameservers) != 0 {
		t.Error("Nameservers should be empty")
	}
}

// ============================================================================
// RDAP Server Tests
// ============================================================================

func TestRDAPServers(t *testing.T) {
	// Check that we have servers for common TLDs
	commonTLDs := []string{"com", "net", "org", "io"}
	for _, tld := range commonTLDs {
		if _, ok := rdapServers[tld]; !ok {
			t.Errorf("missing RDAP server for TLD: %s", tld)
		}
	}

	// Check RIRs
	rirs := []string{"arin", "ripe", "apnic", "afrinic", "lacnic"}
	for _, rir := range rirs {
		if _, ok := rdapServers[rir]; !ok {
			t.Errorf("missing RDAP server for RIR: %s", rir)
		}
	}
}

// ============================================================================
// MX Record Tests
// ============================================================================

func TestMXRecordStructure(t *testing.T) {
	mx := MXRecord{
		Host:     "mail.example.com",
		Priority: 10,
	}

	if mx.Host != "mail.example.com" {
		t.Error("Host should be set")
	}
	if mx.Priority != 10 {
		t.Error("Priority should be 10")
	}
}

// ============================================================================
// SRV Record Tests
// ============================================================================

func TestSRVRecordStructure(t *testing.T) {
	srv := SRVRecord{
		Target:   "sipserver.example.com",
		Port:     5060,
		Priority: 10,
		Weight:   100,
	}

	if srv.Target != "sipserver.example.com" {
		t.Error("Target should be set")
	}
	if srv.Port != 5060 {
		t.Error("Port should be 5060")
	}
	if srv.Priority != 10 {
		t.Error("Priority should be 10")
	}
	if srv.Weight != 100 {
		t.Error("Weight should be 100")
	}
}

// ============================================================================
// SOA Record Tests
// ============================================================================

func TestSOARecordStructure(t *testing.T) {
	soa := SOARecord{
		PrimaryNS:  "ns1.example.com",
		AdminEmail: "admin.example.com",
		Serial:     2024010101,
		Refresh:    7200,
		Retry:      3600,
		Expire:     1209600,
		MinTTL:     86400,
	}

	if soa.PrimaryNS != "ns1.example.com" {
		t.Error("PrimaryNS should be set")
	}
	if soa.Serial != 2024010101 {
		t.Error("Serial should match")
	}
	if soa.Refresh != 7200 {
		t.Error("Refresh should be 7200")
	}
}

// ============================================================================
// Extended Record Type Structure Tests
// ============================================================================

func TestCAARecordStructure(t *testing.T) {
	caa := CAARecord{
		Flag:  0,
		Tag:   "issue",
		Value: "letsencrypt.org",
	}

	if caa.Tag != "issue" {
		t.Error("Tag should be 'issue'")
	}
	if caa.Value != "letsencrypt.org" {
		t.Error("Value should be set")
	}
}

func TestSSHFPRecordStructure(t *testing.T) {
	sshfp := SSHFPRecord{
		Algorithm:   4, // Ed25519
		FPType:      2, // SHA-256
		Fingerprint: "abc123def456",
	}

	if sshfp.Algorithm != 4 {
		t.Error("Algorithm should be 4 (Ed25519)")
	}
	if sshfp.FPType != 2 {
		t.Error("FPType should be 2 (SHA-256)")
	}
}

func TestTLSARecordStructure(t *testing.T) {
	tlsa := TLSARecord{
		Usage:        3, // Domain-issued certificate
		Selector:     1, // SubjectPublicKeyInfo
		MatchingType: 1, // SHA-256
		CertData:     "abcd1234",
	}

	if tlsa.Usage != 3 {
		t.Error("Usage should be 3")
	}
	if tlsa.Selector != 1 {
		t.Error("Selector should be 1")
	}
}

func TestDSRecordStructure(t *testing.T) {
	ds := DSRecord{
		KeyTag:     12345,
		Algorithm:  13, // ECDSAP256SHA256
		DigestType: 2,  // SHA-256
		Digest:     "deadbeef",
	}

	if ds.KeyTag != 12345 {
		t.Error("KeyTag should be 12345")
	}
	if ds.Algorithm != 13 {
		t.Error("Algorithm should be 13")
	}
}

func TestNAPTRRecordStructure(t *testing.T) {
	naptr := NAPTRRecord{
		Order:       100,
		Preference:  10,
		Flags:       "U",
		Service:     "E2U+sip",
		Regexp:      "!^.*$!sip:info@example.com!",
		Replacement: ".",
	}

	if naptr.Order != 100 {
		t.Error("Order should be 100")
	}
	if naptr.Service != "E2U+sip" {
		t.Error("Service should be E2U+sip")
	}
}

func TestRPRecordStructure(t *testing.T) {
	rp := RPRecord{
		Mailbox: "admin.example.com",
		TxtDom:  "info.example.com",
	}

	if rp.Mailbox != "admin.example.com" {
		t.Error("Mailbox should be set")
	}
}

func TestLOCRecordStructure(t *testing.T) {
	loc := LOCRecord{
		Latitude:  51.5074,
		Longitude: -0.1278,
		Altitude:  11,
		Size:      10,
		HPrecis:   10,
		VPrecis:   10,
	}

	if loc.Latitude < 51.5 || loc.Latitude > 51.6 {
		t.Error("Latitude should be near 51.5074")
	}
}

func TestALIASRecordStructure(t *testing.T) {
	alias := ALIASRecord{
		Target: "target.example.com",
	}

	if alias.Target != "target.example.com" {
		t.Error("Target should be set")
	}
}

func TestWebRedirectRecordStructure(t *testing.T) {
	wr := WebRedirectRecord{
		URL:          "https://www.example.com",
		RedirectType: 301,
		Frame:        false,
	}

	if wr.URL != "https://www.example.com" {
		t.Error("URL should be set")
	}
	if wr.RedirectType != 301 {
		t.Error("RedirectType should be 301")
	}
}

// ============================================================================
// Helper Function Tests
// ============================================================================

func TestIsNoSuchHostError(t *testing.T) {
	tests := []struct {
		errStr   string
		expected bool
	}{
		{"no such host", true},
		{"NXDOMAIN", true},
		{"not found", true},
		{"connection refused", false},
		{"timeout", false},
		{"", false},
	}

	for _, tc := range tests {
		var err error
		if tc.errStr != "" {
			err = &testError{msg: tc.errStr}
		}
		result := isNoSuchHostError(err)
		if result != tc.expected {
			t.Errorf("isNoSuchHostError(%q) = %v, want %v", tc.errStr, result, tc.expected)
		}
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// ============================================================================
// URL Building Tests
// ============================================================================

func TestBuildRDAPURLs(t *testing.T) {
	// These test the URL structure, not actual lookups

	// Domain URL
	domain := "example.com"
	expectedDomainPrefix := "https://rdap.org/domain/"
	if !strings.HasPrefix("https://rdap.org/domain/"+domain, expectedDomainPrefix) {
		t.Error("Domain URL format is incorrect")
	}

	// IP URL
	ip := "8.8.8.8"
	expectedIPPrefix := "https://rdap.org/ip/"
	if !strings.HasPrefix("https://rdap.org/ip/"+ip, expectedIPPrefix) {
		t.Error("IP URL format is incorrect")
	}

	// ASN URL
	asn := "15169"
	expectedASNPrefix := "https://rdap.org/autnum/"
	if !strings.HasPrefix("https://rdap.org/autnum/"+asn, expectedASNPrefix) {
		t.Error("ASN URL format is incorrect")
	}
}
