package networktools

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func PrepareProxyHeaders(req *http.Request, target *url.URL) {
	req.Host = target.Host
	req.Header.Set("X-Forwarded-Host", req.Host)

	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	req.Header.Set("X-Real-IP", clientIP)

	if existing := req.Header.Get("X-Forwarded-For"); existing != "" {
		req.Header.Set("X-Forwarded-For", existing+", "+clientIP)
	} else {
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	scheme := "http"
	if IsRequestSecure(req) {
		scheme = "https"
	}
	req.Header.Set("X-Forwarded-Proto", scheme)
}

func GetRealIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func IsRequestSecure(r *http.Request) bool {
	protoHeader := r.Header.Get("X-Forwarded-Proto")
	protoParts := strings.Split(protoHeader, ",")
	proto := ""
	if len(protoParts) > 0 {
		proto = strings.ToLower(strings.TrimSpace(protoParts[0]))
	}
	return r.TLS != nil || proto == "https"
}

func BuildHostURL(schema, host, port string, paths ...string) string {
	host = strings.TrimRight(host, "/")
	parsed, err := url.Parse(host)
	if err == nil && parsed.Host != "" {
		host = parsed.Host
	}

	cleanHost, _, err := net.SplitHostPort(host)
	if err != nil {
		cleanHost = host
	}

	if port != "" && port != "80" && port != "443" {
		cleanHost = fmt.Sprintf("%s:%s", cleanHost, port)
	}

	fullPath := ""
	for _, p := range paths {
		fullPath += "/" + strings.Trim(p, "/")
	}

	return fmt.Sprintf("%s://%s%s", schema, cleanHost, fullPath)
}
