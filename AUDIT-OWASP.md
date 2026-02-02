# OWASP Top 10 Security Audit

## Summary
0 critical, 2 high, 1 medium findings

## Findings by Category

### A06: Vulnerable and Outdated Components (High)
- **Finding:** The `govulncheck` tool identified 13 vulnerabilities in the Go standard library, stemming from an outdated Go version.
- **Remediation:** It is recommended to upgrade the project's Go version to the latest stable release to mitigate these vulnerabilities.

### A10: Server-Side Request Forgery (SSRF) (High)
- **Finding:** The `RDAPLookupDomainWithTimeout`, `RDAPLookupIPWithTimeout`, and `RDAPLookupASNWithTimeout` functions constructed request URLs by directly embedding user-provided inputs. This could have allowed a malicious actor to craft inputs that would cause the server to make requests to internal resources.
- **Remediation:** All user-provided inputs (`domain`, `ip`, and `asn`) are now sanitized using `url.PathEscape()` before being included in the request URL, preventing path traversal and other SSRF-style attacks.

### A03: Injection (Medium)
- **Finding:** The `DNSLookup...` functions did not sanitize the `domain` parameter, which could have led to unexpected behavior if special characters were provided as input.
- **Remediation:** The `domain` parameter is now validated using a regular expression to ensure it conforms to a valid domain name format, mitigating the risk of injection attacks.
