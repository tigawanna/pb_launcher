package networktools

import (
	"net"
	"net/http"
	"net/url"
)

func PrepareProxyHeaders(req *http.Request, target *url.URL) {
	req.Host = target.Host
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Header.Set("X-Real-IP", GetRealIP(req))

	clientIP, _, _ := net.SplitHostPort(req.RemoteAddr)
	if existing := req.Header.Get("X-Forwarded-For"); existing != "" {
		req.Header.Set("X-Forwarded-For", existing+", "+clientIP)
	} else {
		req.Header.Set("X-Forwarded-For", clientIP)
	}

	req.Header.Set("X-Forwarded-Proto", Scheme(req))
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

func Scheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}
