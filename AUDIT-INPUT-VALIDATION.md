# Input Validation and Sanitization Audit

## Input Entry Points Inventory

This section inventories all the points where untrusted input enters the system.

### Go API (`poindexter.go`)

- **`Hello(name string)`**: Accepts a string input, which is used in a greeting message.

### WebAssembly Interface (`wasm/main.go`)

The following functions are exposed to JavaScript and represent the primary attack surface from a web environment:

- **`pxHello(name string)`**: Accepts a string from JS.
- **`pxNewTree(dim int)`**: Accepts an integer dimension.
- **`pxInsert(treeId int, point js.Value)`**: Accepts a JS object with `id` (string), `coords` (array of numbers), and `value` (string).
- **`pxDeleteByID(treeId int, id string)`**: Accepts a string point ID.
- **`pxNearest(treeId int, query []float64)`**: Accepts an array of numbers.
- **`pxKNearest(treeId int, query []float64, k int)`**: Accepts an array of numbers and an integer.
- **`pxRadius(treeId int, query []float64, r float64)`**: Accepts an array of numbers and a float.
- **`pxComputePeerQualityScore(metrics js.Value, weights js.Value)`**: Accepts JS objects for metrics and weights, containing various numeric and string fields.
- **`pxComputeTrustScore(metrics js.Value)`**: Accepts a JS object for trust metrics, containing various numeric and string fields.
- **`pxGetExternalToolLinks(domain string)`**: Accepts a domain string.
- **`pxGetExternalToolLinksIP(ip string)`**: Accepts an IP address string.
- **`pxGetExternalToolLinksEmail(emailOrDomain string)`**: Accepts an email or domain string.
- **`pxBuildRDAPDomainURL(domain string)`**: Accepts a domain string.
- **`pxBuildRDAPIPURL(ip string)`**: Accepts an IP address string.
- **`pxBuildRDAPASNURL(asn string)`**: Accepts an ASN string.

### DNS Tools (`dns_tools.go`)

These functions are called internally but are also exposed via the WASM interface. They interact with external DNS resolvers and RDAP servers.

- **`DNSLookup(domain string, recordType DNSRecordType)`**: Accepts a domain string.
- **`RDAPLookupDomain(domain string)`**: Accepts a domain string.
- **`RDAPLookupIP(ip string)`**: Accepts an IP address string.
- **`RDAPLookupASN(asn string)`**: Accepts an ASN string.
- **`GetExternalToolLinks(domain string)`**: Accepts a domain string.
- **`GetExternalToolLinksIP(ip string)`**: Accepts an IP address string.
- **`GetExternalToolLinksEmail(emailOrDomain string)`**: Accepts an email or domain string.

## Validation Gaps Found

- **WASM Interface**: The functions exposed to JavaScript in `wasm/main.go` perform minimal validation. For example, `pxComputePeerQualityScore` and `pxComputeTrustScore` accept numeric and string inputs from JS objects without proper sanitization.
- **DNS Tools**: The `GetExternalToolLinks` functions in `dns_tools.go` use `url.QueryEscape` but do not validate the input strings for correctness, potentially allowing malformed data to be passed to external services.

## Injection Vectors Discovered

- **URL Injection**: The `GetExternalToolLinks` functions could be vulnerable to URL injection if an attacker provides a maliciously crafted domain or IP address. For example, an input like `example.com?some_param=some_value` could alter the behavior of the external service being linked to.

## Remediation Recommendations

To address the identified vulnerabilities, the following remediation actions are recommended:

### 1. Sanitize Inputs in WASM Interface

All inputs from JavaScript should be treated as untrusted and validated before use.

**Example: `pxComputePeerQualityScore`**

```go
// wasm/main.go

import (
	"errors"
	"regexp"
)

var (
	// Basic validation for domain names and IPs
	domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	ipRegex     = regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
)

func computePeerQualityScore(_ js.Value, args []js.Value) (any, error) {
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

	// Add validation for NATType
	if !isValidNATType(metrics.NATType) {
		return nil, errors.New("invalid NATType")
	}

	// ... rest of the function
}

func isValidNATType(natType string) bool {
	// Implement validation logic for NATType
	return natType == "static" || natType == "dynamic"
}
```

### 2. Validate and Sanitize DNS Tool Inputs

Before creating external links, validate that the input is a valid domain, IP, or email.

**Example: `GetExternalToolLinks`**

```go
// dns_tools.go

import (
	"net"
	"net/url"
	"regexp"
)

var (
	domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
)

func GetExternalToolLinks(domain string) ExternalToolLinks {
	if !domainRegex.MatchString(domain) {
		// Return empty links or an error
		return ExternalToolLinks{Target: domain, Type: "domain"}
	}
	encoded := url.QueryEscape(domain)
	// ... rest of the function
}

func GetExternalToolLinksIP(ip string) ExternalToolLinks {
	if net.ParseIP(ip) == nil {
		return ExternalToolLinks{Target: ip, Type: "ip"}
	}
	encoded := url.QueryEscape(ip)
	// ... rest of the function
}
```
