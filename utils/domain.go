package utils

import "strings"

func NormalizeWildcardDomain(domain string) string {
	if !strings.HasPrefix(domain, "*.") {
		return "*." + domain
	}
	return domain
}
