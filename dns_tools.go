package poindexter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// ============================================================================
// DNS Record Types
// ============================================================================

// DNSRecordType represents DNS record types
type DNSRecordType string

const (
	DNSRecordA     DNSRecordType = "A"
	DNSRecordAAAA  DNSRecordType = "AAAA"
	DNSRecordMX    DNSRecordType = "MX"
	DNSRecordTXT   DNSRecordType = "TXT"
	DNSRecordNS    DNSRecordType = "NS"
	DNSRecordCNAME DNSRecordType = "CNAME"
	DNSRecordSOA   DNSRecordType = "SOA"
	DNSRecordPTR   DNSRecordType = "PTR"
	DNSRecordSRV   DNSRecordType = "SRV"
	DNSRecordCAA   DNSRecordType = "CAA"
)

// DNSRecord represents a generic DNS record
type DNSRecord struct {
	Type  DNSRecordType `json:"type"`
	Name  string        `json:"name"`
	Value string        `json:"value"`
	TTL   int           `json:"ttl,omitempty"`
}

// MXRecord represents an MX record with priority
type MXRecord struct {
	Host     string `json:"host"`
	Priority uint16 `json:"priority"`
}

// SRVRecord represents an SRV record
type SRVRecord struct {
	Target   string `json:"target"`
	Port     uint16 `json:"port"`
	Priority uint16 `json:"priority"`
	Weight   uint16 `json:"weight"`
}

// SOARecord represents an SOA record
type SOARecord struct {
	PrimaryNS  string `json:"primaryNs"`
	AdminEmail string `json:"adminEmail"`
	Serial     uint32 `json:"serial"`
	Refresh    uint32 `json:"refresh"`
	Retry      uint32 `json:"retry"`
	Expire     uint32 `json:"expire"`
	MinTTL     uint32 `json:"minTtl"`
}

// DNSLookupResult contains the results of a DNS lookup
type DNSLookupResult struct {
	Domain      string      `json:"domain"`
	QueryType   string      `json:"queryType"`
	Records     []DNSRecord `json:"records"`
	MXRecords   []MXRecord  `json:"mxRecords,omitempty"`
	SRVRecords  []SRVRecord `json:"srvRecords,omitempty"`
	SOARecord   *SOARecord  `json:"soaRecord,omitempty"`
	LookupTimeMs int64      `json:"lookupTimeMs"`
	Error       string      `json:"error,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}

// CompleteDNSLookup contains all DNS records for a domain
type CompleteDNSLookup struct {
	Domain       string          `json:"domain"`
	A            []string        `json:"a,omitempty"`
	AAAA         []string        `json:"aaaa,omitempty"`
	MX           []MXRecord      `json:"mx,omitempty"`
	NS           []string        `json:"ns,omitempty"`
	TXT          []string        `json:"txt,omitempty"`
	CNAME        string          `json:"cname,omitempty"`
	SOA          *SOARecord      `json:"soa,omitempty"`
	LookupTimeMs int64           `json:"lookupTimeMs"`
	Errors       []string        `json:"errors,omitempty"`
	Timestamp    time.Time       `json:"timestamp"`
}

// ============================================================================
// DNS Lookup Functions
// ============================================================================

// DNSLookup performs a DNS lookup for the specified record type
func DNSLookup(domain string, recordType DNSRecordType) DNSLookupResult {
	return DNSLookupWithTimeout(domain, recordType, 10*time.Second)
}

// DNSLookupWithTimeout performs a DNS lookup with a custom timeout
func DNSLookupWithTimeout(domain string, recordType DNSRecordType, timeout time.Duration) DNSLookupResult {
	start := time.Now()
	result := DNSLookupResult{
		Domain:    domain,
		QueryType: string(recordType),
		Timestamp: start,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver := net.Resolver{}

	switch recordType {
	case DNSRecordA:
		ips, err := resolver.LookupIP(ctx, "ip4", domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, ip := range ips {
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordA,
					Name:  domain,
					Value: ip.String(),
				})
			}
		}

	case DNSRecordAAAA:
		ips, err := resolver.LookupIP(ctx, "ip6", domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, ip := range ips {
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordAAAA,
					Name:  domain,
					Value: ip.String(),
				})
			}
		}

	case DNSRecordMX:
		mxs, err := resolver.LookupMX(ctx, domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, mx := range mxs {
				result.MXRecords = append(result.MXRecords, MXRecord{
					Host:     strings.TrimSuffix(mx.Host, "."),
					Priority: mx.Pref,
				})
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordMX,
					Name:  domain,
					Value: fmt.Sprintf("%d %s", mx.Pref, mx.Host),
				})
			}
			// Sort by priority
			sort.Slice(result.MXRecords, func(i, j int) bool {
				return result.MXRecords[i].Priority < result.MXRecords[j].Priority
			})
		}

	case DNSRecordTXT:
		txts, err := resolver.LookupTXT(ctx, domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, txt := range txts {
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordTXT,
					Name:  domain,
					Value: txt,
				})
			}
		}

	case DNSRecordNS:
		nss, err := resolver.LookupNS(ctx, domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, ns := range nss {
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordNS,
					Name:  domain,
					Value: strings.TrimSuffix(ns.Host, "."),
				})
			}
		}

	case DNSRecordCNAME:
		cname, err := resolver.LookupCNAME(ctx, domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Records = append(result.Records, DNSRecord{
				Type:  DNSRecordCNAME,
				Name:  domain,
				Value: strings.TrimSuffix(cname, "."),
			})
		}

	case DNSRecordSRV:
		// SRV records require a service and protocol prefix, e.g., _http._tcp.example.com
		_, srvs, err := resolver.LookupSRV(ctx, "", "", domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, srv := range srvs {
				result.SRVRecords = append(result.SRVRecords, SRVRecord{
					Target:   strings.TrimSuffix(srv.Target, "."),
					Port:     srv.Port,
					Priority: srv.Priority,
					Weight:   srv.Weight,
				})
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordSRV,
					Name:  domain,
					Value: fmt.Sprintf("%d %d %d %s", srv.Priority, srv.Weight, srv.Port, srv.Target),
				})
			}
		}

	case DNSRecordPTR:
		names, err := resolver.LookupAddr(ctx, domain)
		if err != nil {
			result.Error = err.Error()
		} else {
			for _, name := range names {
				result.Records = append(result.Records, DNSRecord{
					Type:  DNSRecordPTR,
					Name:  domain,
					Value: strings.TrimSuffix(name, "."),
				})
			}
		}

	default:
		result.Error = fmt.Sprintf("unsupported record type: %s", recordType)
	}

	result.LookupTimeMs = time.Since(start).Milliseconds()
	return result
}

// DNSLookupAll performs lookups for all common record types
func DNSLookupAll(domain string) CompleteDNSLookup {
	return DNSLookupAllWithTimeout(domain, 10*time.Second)
}

// DNSLookupAllWithTimeout performs lookups for all common record types with timeout
func DNSLookupAllWithTimeout(domain string, timeout time.Duration) CompleteDNSLookup {
	start := time.Now()
	result := CompleteDNSLookup{
		Domain:    domain,
		Timestamp: start,
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver := net.Resolver{}

	// A records
	if ips, err := resolver.LookupIP(ctx, "ip4", domain); err == nil {
		for _, ip := range ips {
			result.A = append(result.A, ip.String())
		}
	} else if !isNoSuchHostError(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("A: %s", err.Error()))
	}

	// AAAA records
	if ips, err := resolver.LookupIP(ctx, "ip6", domain); err == nil {
		for _, ip := range ips {
			result.AAAA = append(result.AAAA, ip.String())
		}
	} else if !isNoSuchHostError(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("AAAA: %s", err.Error()))
	}

	// MX records
	if mxs, err := resolver.LookupMX(ctx, domain); err == nil {
		for _, mx := range mxs {
			result.MX = append(result.MX, MXRecord{
				Host:     strings.TrimSuffix(mx.Host, "."),
				Priority: mx.Pref,
			})
		}
		sort.Slice(result.MX, func(i, j int) bool {
			return result.MX[i].Priority < result.MX[j].Priority
		})
	} else if !isNoSuchHostError(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("MX: %s", err.Error()))
	}

	// NS records
	if nss, err := resolver.LookupNS(ctx, domain); err == nil {
		for _, ns := range nss {
			result.NS = append(result.NS, strings.TrimSuffix(ns.Host, "."))
		}
	} else if !isNoSuchHostError(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("NS: %s", err.Error()))
	}

	// TXT records
	if txts, err := resolver.LookupTXT(ctx, domain); err == nil {
		result.TXT = txts
	} else if !isNoSuchHostError(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("TXT: %s", err.Error()))
	}

	// CNAME record
	if cname, err := resolver.LookupCNAME(ctx, domain); err == nil {
		result.CNAME = strings.TrimSuffix(cname, ".")
		// If CNAME equals domain, it's not really a CNAME
		if result.CNAME == domain {
			result.CNAME = ""
		}
	}

	result.LookupTimeMs = time.Since(start).Milliseconds()
	return result
}

// ReverseDNSLookup performs a reverse DNS lookup for an IP address
func ReverseDNSLookup(ip string) DNSLookupResult {
	return DNSLookupWithTimeout(ip, DNSRecordPTR, 10*time.Second)
}

func isNoSuchHostError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no such host") ||
		strings.Contains(err.Error(), "NXDOMAIN") ||
		strings.Contains(err.Error(), "not found")
}

// ============================================================================
// RDAP (Registration Data Access Protocol) - New Style WHOIS
// ============================================================================

// RDAPResponse represents an RDAP response
type RDAPResponse struct {
	// Common fields
	Handle      string       `json:"handle,omitempty"`
	LDHName     string       `json:"ldhName,omitempty"` // Domain name
	UnicodeName string       `json:"unicodeName,omitempty"`
	Status      []string     `json:"status,omitempty"`
	Events      []RDAPEvent  `json:"events,omitempty"`
	Entities    []RDAPEntity `json:"entities,omitempty"`
	Nameservers []RDAPNs     `json:"nameservers,omitempty"`
	Links       []RDAPLink   `json:"links,omitempty"`
	Remarks     []RDAPRemark `json:"remarks,omitempty"`
	Notices     []RDAPNotice `json:"notices,omitempty"`

	// Network-specific (for IP lookups)
	StartAddress string `json:"startAddress,omitempty"`
	EndAddress   string `json:"endAddress,omitempty"`
	IPVersion    string `json:"ipVersion,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Country      string `json:"country,omitempty"`
	ParentHandle string `json:"parentHandle,omitempty"`

	// Error fields
	ErrorCode    int    `json:"errorCode,omitempty"`
	Title        string `json:"title,omitempty"`
	Description  []string `json:"description,omitempty"`

	// Metadata
	RawJSON      string    `json:"rawJson,omitempty"`
	LookupTimeMs int64     `json:"lookupTimeMs"`
	Timestamp    time.Time `json:"timestamp"`
	Error        string    `json:"error,omitempty"`
}

// RDAPEvent represents an RDAP event (registration, expiration, etc.)
type RDAPEvent struct {
	EventAction string `json:"eventAction"`
	EventDate   string `json:"eventDate"`
	EventActor  string `json:"eventActor,omitempty"`
}

// RDAPEntity represents an entity (registrar, registrant, etc.)
type RDAPEntity struct {
	Handle     string       `json:"handle,omitempty"`
	Roles      []string     `json:"roles,omitempty"`
	VCardArray []any        `json:"vcardArray,omitempty"`
	Entities   []RDAPEntity `json:"entities,omitempty"`
	Events     []RDAPEvent  `json:"events,omitempty"`
	Links      []RDAPLink   `json:"links,omitempty"`
	Remarks    []RDAPRemark `json:"remarks,omitempty"`
}

// RDAPNs represents a nameserver in RDAP
type RDAPNs struct {
	LDHName     string   `json:"ldhName"`
	IPAddresses *RDAPIPs `json:"ipAddresses,omitempty"`
}

// RDAPIPs represents IP addresses for a nameserver
type RDAPIPs struct {
	V4 []string `json:"v4,omitempty"`
	V6 []string `json:"v6,omitempty"`
}

// RDAPLink represents a link in RDAP
type RDAPLink struct {
	Value string `json:"value,omitempty"`
	Rel   string `json:"rel,omitempty"`
	Href  string `json:"href,omitempty"`
	Type  string `json:"type,omitempty"`
}

// RDAPRemark represents a remark/notice
type RDAPRemark struct {
	Title       string   `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Links       []RDAPLink `json:"links,omitempty"`
}

// RDAPNotice is an alias for RDAPRemark
type RDAPNotice = RDAPRemark

// RDAPBootstrapRegistry holds the RDAP bootstrap data
type RDAPBootstrapRegistry struct {
	Services [][]interface{} `json:"services"`
	Version  string          `json:"version"`
}

// RDAP server URLs for different TLDs and RIRs
var rdapServers = map[string]string{
	// Generic TLDs (ICANN)
	"com":  "https://rdap.verisign.com/com/v1/",
	"net":  "https://rdap.verisign.com/net/v1/",
	"org":  "https://rdap.publicinterestregistry.org/rdap/",
	"info": "https://rdap.afilias.net/rdap/info/",
	"biz":  "https://rdap.afilias.net/rdap/biz/",
	"io":   "https://rdap.nic.io/",
	"co":   "https://rdap.nic.co/",
	"me":   "https://rdap.nic.me/",
	"app":  "https://rdap.nic.google/",
	"dev":  "https://rdap.nic.google/",

	// Country code TLDs
	"uk":   "https://rdap.nominet.uk/uk/",
	"de":   "https://rdap.denic.de/",
	"nl":   "https://rdap.sidn.nl/",
	"au":   "https://rdap.auda.org.au/",
	"nz":   "https://rdap.dns.net.nz/",
	"br":   "https://rdap.registro.br/",
	"jp":   "https://rdap.jprs.jp/",

	// RIRs for IP lookups
	"arin":    "https://rdap.arin.net/registry/",
	"ripe":    "https://rdap.db.ripe.net/",
	"apnic":   "https://rdap.apnic.net/",
	"afrinic": "https://rdap.afrinic.net/rdap/",
	"lacnic":  "https://rdap.lacnic.net/rdap/",
}

// RDAPLookupDomain performs an RDAP lookup for a domain
func RDAPLookupDomain(domain string) RDAPResponse {
	return RDAPLookupDomainWithTimeout(domain, 15*time.Second)
}

// RDAPLookupDomainWithTimeout performs an RDAP lookup with custom timeout
func RDAPLookupDomainWithTimeout(domain string, timeout time.Duration) RDAPResponse {
	start := time.Now()
	result := RDAPResponse{
		LDHName:   domain,
		Timestamp: start,
	}

	// Extract TLD
	parts := strings.Split(strings.ToLower(domain), ".")
	if len(parts) < 2 {
		result.Error = "invalid domain format"
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}
	tld := parts[len(parts)-1]

	// Find RDAP server
	serverURL, ok := rdapServers[tld]
	if !ok {
		// Try to use IANA bootstrap
		serverURL = fmt.Sprintf("https://rdap.org/domain/%s", domain)
	} else {
		serverURL = serverURL + "domain/" + domain
	}

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(serverURL)
	if err != nil {
		result.Error = fmt.Sprintf("RDAP request failed: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read response: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	result.RawJSON = string(body)

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("RDAP server returned status %d", resp.StatusCode)
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		result.Error = fmt.Sprintf("failed to parse RDAP response: %s", err.Error())
	}

	result.LookupTimeMs = time.Since(start).Milliseconds()
	return result
}

// RDAPLookupIP performs an RDAP lookup for an IP address
func RDAPLookupIP(ip string) RDAPResponse {
	return RDAPLookupIPWithTimeout(ip, 15*time.Second)
}

// RDAPLookupIPWithTimeout performs an RDAP lookup for an IP with custom timeout
func RDAPLookupIPWithTimeout(ip string, timeout time.Duration) RDAPResponse {
	start := time.Now()
	result := RDAPResponse{
		StartAddress: ip,
		Timestamp:    start,
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		result.Error = "invalid IP address"
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	// Use rdap.org as a universal redirector
	serverURL := fmt.Sprintf("https://rdap.org/ip/%s", ip)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(serverURL)
	if err != nil {
		result.Error = fmt.Sprintf("RDAP request failed: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read response: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	result.RawJSON = string(body)

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("RDAP server returned status %d", resp.StatusCode)
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		result.Error = fmt.Sprintf("failed to parse RDAP response: %s", err.Error())
	}

	result.LookupTimeMs = time.Since(start).Milliseconds()
	return result
}

// RDAPLookupASN performs an RDAP lookup for an ASN
func RDAPLookupASN(asn string) RDAPResponse {
	return RDAPLookupASNWithTimeout(asn, 15*time.Second)
}

// RDAPLookupASNWithTimeout performs an RDAP lookup for an ASN with timeout
func RDAPLookupASNWithTimeout(asn string, timeout time.Duration) RDAPResponse {
	start := time.Now()
	result := RDAPResponse{
		Handle:    asn,
		Timestamp: start,
	}

	// Normalize ASN (remove "AS" prefix if present)
	asnNum := strings.TrimPrefix(strings.ToUpper(asn), "AS")

	// Use rdap.org as a universal redirector
	serverURL := fmt.Sprintf("https://rdap.org/autnum/%s", asnNum)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Get(serverURL)
	if err != nil {
		result.Error = fmt.Sprintf("RDAP request failed: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read response: %s", err.Error())
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	result.RawJSON = string(body)

	if resp.StatusCode != http.StatusOK {
		result.Error = fmt.Sprintf("RDAP server returned status %d", resp.StatusCode)
		result.LookupTimeMs = time.Since(start).Milliseconds()
		return result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		result.Error = fmt.Sprintf("failed to parse RDAP response: %s", err.Error())
	}

	result.LookupTimeMs = time.Since(start).Milliseconds()
	return result
}

// ============================================================================
// External Tool Links
// ============================================================================

// ExternalToolLinks contains links to external DNS/network analysis tools
type ExternalToolLinks struct {
	// Target being analyzed
	Target string `json:"target"`
	Type   string `json:"type"` // "domain", "ip", "email"

	// MXToolbox links
	MXToolboxDNS         string `json:"mxtoolboxDns,omitempty"`
	MXToolboxMX          string `json:"mxtoolboxMx,omitempty"`
	MXToolboxBlacklist   string `json:"mxtoolboxBlacklist,omitempty"`
	MXToolboxSMTP        string `json:"mxtoolboxSmtp,omitempty"`
	MXToolboxSPF         string `json:"mxtoolboxSpf,omitempty"`
	MXToolboxDMARC       string `json:"mxtoolboxDmarc,omitempty"`
	MXToolboxDKIM        string `json:"mxtoolboxDkim,omitempty"`
	MXToolboxHTTP        string `json:"mxtoolboxHttp,omitempty"`
	MXToolboxHTTPS       string `json:"mxtoolboxHttps,omitempty"`
	MXToolboxPing        string `json:"mxtoolboxPing,omitempty"`
	MXToolboxTrace       string `json:"mxtoolboxTrace,omitempty"`
	MXToolboxWhois       string `json:"mxtoolboxWhois,omitempty"`
	MXToolboxASN         string `json:"mxtoolboxAsn,omitempty"`

	// DNSChecker links
	DNSCheckerDNS        string `json:"dnscheckerDns,omitempty"`
	DNSCheckerPropagation string `json:"dnscheckerPropagation,omitempty"`

	// Other tools
	WhoIs                string `json:"whois,omitempty"`
	ViewDNS              string `json:"viewdns,omitempty"`
	IntoDNS              string `json:"intodns,omitempty"`
	DNSViz               string `json:"dnsviz,omitempty"`
	SecurityTrails       string `json:"securitytrails,omitempty"`
	Shodan               string `json:"shodan,omitempty"`
	Censys               string `json:"censys,omitempty"`
	BuiltWith            string `json:"builtwith,omitempty"`
	SSLLabs              string `json:"ssllabs,omitempty"`
	HSTSPreload          string `json:"hstsPreload,omitempty"`
	Hardenize            string `json:"hardenize,omitempty"`

	// IP-specific tools
	IPInfo               string `json:"ipinfo,omitempty"`
	AbuseIPDB            string `json:"abuseipdb,omitempty"`
	VirusTotal           string `json:"virustotal,omitempty"`
	ThreatCrowd          string `json:"threatcrowd,omitempty"`

	// Email-specific tools
	MailTester           string `json:"mailtester,omitempty"`
	LearnDMARC           string `json:"learndmarc,omitempty"`
}

// GetExternalToolLinks generates links to external analysis tools for a domain
func GetExternalToolLinks(domain string) ExternalToolLinks {
	encoded := url.QueryEscape(domain)

	return ExternalToolLinks{
		Target: domain,
		Type:   "domain",

		// MXToolbox
		MXToolboxDNS:       fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=dns%%3a%s&run=toolpage", encoded),
		MXToolboxMX:        fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=mx%%3a%s&run=toolpage", encoded),
		MXToolboxBlacklist: fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=blacklist%%3a%s&run=toolpage", encoded),
		MXToolboxSMTP:      fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=smtp%%3a%s&run=toolpage", encoded),
		MXToolboxSPF:       fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=spf%%3a%s&run=toolpage", encoded),
		MXToolboxDMARC:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=dmarc%%3a%s&run=toolpage", encoded),
		MXToolboxDKIM:      fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=dkim%%3a%s&run=toolpage", encoded),
		MXToolboxHTTP:      fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=http%%3a%s&run=toolpage", encoded),
		MXToolboxHTTPS:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=https%%3a%s&run=toolpage", encoded),
		MXToolboxPing:      fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=ping%%3a%s&run=toolpage", encoded),
		MXToolboxTrace:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=trace%%3a%s&run=toolpage", encoded),
		MXToolboxWhois:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=whois%%3a%s&run=toolpage", encoded),

		// DNSChecker
		DNSCheckerDNS:        fmt.Sprintf("https://dnschecker.org/#A/%s", encoded),
		DNSCheckerPropagation: fmt.Sprintf("https://dnschecker.org/dns-propagation.php?domain=%s", encoded),

		// Other tools
		WhoIs:          fmt.Sprintf("https://who.is/whois/%s", encoded),
		ViewDNS:        fmt.Sprintf("https://viewdns.info/dnsrecord/?domain=%s", encoded),
		IntoDNS:        fmt.Sprintf("https://intodns.com/%s", encoded),
		DNSViz:         fmt.Sprintf("https://dnsviz.net/d/%s/analyze/", encoded),
		SecurityTrails: fmt.Sprintf("https://securitytrails.com/domain/%s", encoded),
		BuiltWith:      fmt.Sprintf("https://builtwith.com/%s", encoded),
		SSLLabs:        fmt.Sprintf("https://www.ssllabs.com/ssltest/analyze.html?d=%s", encoded),
		HSTSPreload:    fmt.Sprintf("https://hstspreload.org/?domain=%s", encoded),
		Hardenize:      fmt.Sprintf("https://www.hardenize.com/report/%s", encoded),
		VirusTotal:     fmt.Sprintf("https://www.virustotal.com/gui/domain/%s", encoded),
	}
}

// GetExternalToolLinksIP generates links to external analysis tools for an IP
func GetExternalToolLinksIP(ip string) ExternalToolLinks {
	encoded := url.QueryEscape(ip)

	return ExternalToolLinks{
		Target: ip,
		Type:   "ip",

		// MXToolbox
		MXToolboxBlacklist: fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=blacklist%%3a%s&run=toolpage", encoded),
		MXToolboxPing:      fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=ping%%3a%s&run=toolpage", encoded),
		MXToolboxTrace:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=trace%%3a%s&run=toolpage", encoded),
		MXToolboxWhois:     fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=whois%%3a%s&run=toolpage", encoded),
		MXToolboxASN:       fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=asn%%3a%s&run=toolpage", encoded),

		// IP-specific tools
		IPInfo:      fmt.Sprintf("https://ipinfo.io/%s", encoded),
		AbuseIPDB:   fmt.Sprintf("https://www.abuseipdb.com/check/%s", encoded),
		VirusTotal:  fmt.Sprintf("https://www.virustotal.com/gui/ip-address/%s", encoded),
		Shodan:      fmt.Sprintf("https://www.shodan.io/host/%s", encoded),
		Censys:      fmt.Sprintf("https://search.censys.io/hosts/%s", encoded),
		ThreatCrowd: fmt.Sprintf("https://www.threatcrowd.org/ip.php?ip=%s", encoded),
	}
}

// GetExternalToolLinksEmail generates links for email-related checks
func GetExternalToolLinksEmail(emailOrDomain string) ExternalToolLinks {
	// Extract domain from email if needed
	domain := emailOrDomain
	if strings.Contains(emailOrDomain, "@") {
		parts := strings.Split(emailOrDomain, "@")
		if len(parts) == 2 {
			domain = parts[1]
		}
	}

	encoded := url.QueryEscape(domain)
	emailEncoded := url.QueryEscape(emailOrDomain)

	return ExternalToolLinks{
		Target: emailOrDomain,
		Type:   "email",

		// MXToolbox email checks
		MXToolboxMX:    fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=mx%%3a%s&run=toolpage", encoded),
		MXToolboxSMTP:  fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=smtp%%3a%s&run=toolpage", encoded),
		MXToolboxSPF:   fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=spf%%3a%s&run=toolpage", encoded),
		MXToolboxDMARC: fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=dmarc%%3a%s&run=toolpage", encoded),
		MXToolboxDKIM:  fmt.Sprintf("https://mxtoolbox.com/SuperTool.aspx?action=dkim%%3a%s&run=toolpage", encoded),

		// Email-specific tools
		MailTester: fmt.Sprintf("https://www.mail-tester.com/test-%s", emailEncoded),
		LearnDMARC: fmt.Sprintf("https://www.learndmarc.com/?domain=%s", encoded),
	}
}

// ============================================================================
// Convenience Types for Parsed Results
// ============================================================================

// ParsedDomainInfo provides a simplified view of domain information
type ParsedDomainInfo struct {
	Domain           string    `json:"domain"`
	Registrar        string    `json:"registrar,omitempty"`
	RegistrationDate string    `json:"registrationDate,omitempty"`
	ExpirationDate   string    `json:"expirationDate,omitempty"`
	UpdatedDate      string    `json:"updatedDate,omitempty"`
	Status           []string  `json:"status,omitempty"`
	Nameservers      []string  `json:"nameservers,omitempty"`
	DNSSEC           bool      `json:"dnssec"`
}

// ParseRDAPResponse extracts key information from an RDAP response
func ParseRDAPResponse(resp RDAPResponse) ParsedDomainInfo {
	info := ParsedDomainInfo{
		Domain:      resp.LDHName,
		Status:      resp.Status,
	}

	// Extract dates from events
	for _, event := range resp.Events {
		switch event.EventAction {
		case "registration":
			info.RegistrationDate = event.EventDate
		case "expiration":
			info.ExpirationDate = event.EventDate
		case "last changed", "last update":
			info.UpdatedDate = event.EventDate
		}
	}

	// Extract registrar from entities
	for _, entity := range resp.Entities {
		for _, role := range entity.Roles {
			if role == "registrar" {
				info.Registrar = entity.Handle
				break
			}
		}
	}

	// Extract nameservers
	for _, ns := range resp.Nameservers {
		info.Nameservers = append(info.Nameservers, ns.LDHName)
	}

	// Check for DNSSEC
	for _, status := range resp.Status {
		if strings.Contains(strings.ToLower(status), "dnssec") ||
			strings.Contains(strings.ToLower(status), "signed") {
			info.DNSSEC = true
			break
		}
	}

	return info
}
