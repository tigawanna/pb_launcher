package domainutil

import "strings"

func ToWildcardDomain(domain string) string {
	domain = strings.TrimPrefix(domain, "*.")
	return "*." + domain
}

func BaseDomain(domain string) string {
	if after, ok := strings.CutPrefix(domain, "*."); ok {
		return after
	}
	return domain
}

func IsWildcardDomain(domain string) bool {
	return strings.HasPrefix(domain, "*.")
}

func SubdomainMatchesWildcard(subdomain, wildcard string) bool {
	if !IsWildcardDomain(wildcard) {
		return false
	}
	wildcardBase := strings.TrimPrefix(wildcard, "*.")
	return strings.HasSuffix(subdomain, "."+wildcardBase) || subdomain == wildcardBase
}
