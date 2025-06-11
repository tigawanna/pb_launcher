package proxy

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"pb_launcher/configs"
	"pb_launcher/internal/proxy/domain"
	"pb_launcher/utils/networktools"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/apis"
)

type DynamicReverseProxy struct {
	discovery       *domain.ServiceDiscovery
	domainDiscovery *domain.DomainServiceDiscovery
	apiDomain       string
	apiAddress      string
	useHttps        bool
	httpsPort       string
	timeout         time.Duration
}

var _ http.Handler = (*DynamicReverseProxy)(nil)

func NewDynamicReverseProxy(
	discovery *domain.ServiceDiscovery,
	domainDiscovery *domain.DomainServiceDiscovery,
	cfg configs.Config,
	pbConf *apis.ServeConfig,
) *DynamicReverseProxy {
	return &DynamicReverseProxy{
		discovery:       discovery,
		domainDiscovery: domainDiscovery,
		apiDomain:       cfg.GetDomain(),
		useHttps:        cfg.UseHttps(),
		httpsPort:       cfg.GetBindHttpsPort(),
		apiAddress:      pbConf.HttpAddr,
		timeout:         15 * time.Second,
	}
}

func (rp *DynamicReverseProxy) resolveServiceID(ctx context.Context, host string) (string, error) {
	if host == rp.apiDomain {
		return "", nil
	}

	if strings.HasSuffix(host, "."+rp.apiDomain) {
		prefix := strings.TrimSuffix(host, "."+rp.apiDomain)
		if strings.Contains(prefix, ".") || prefix == "" {
			return "", fmt.Errorf("invalid service ID")
		}
		return prefix, nil
	}

	serviceID, err := rp.domainDiscovery.FindServiceIDByDomain(ctx, host)
	if err != nil {
		return "", errors.New("service ID not found for domain")
	}
	return *serviceID, nil
}

func (rp *DynamicReverseProxy) resolveServiceTarget(ctx context.Context, host string) (*url.URL, error) {
	serviceID, err := rp.resolveServiceID(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("invalid service ID: %w", err)
	}

	if serviceID == "" {
		return &url.URL{
			Scheme: "http",
			Host:   rp.apiAddress,
		}, nil
	}
	service, err := rp.discovery.FindRunningServiceByID(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	return &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(service.IP, strconv.Itoa(service.Port)),
	}, nil
}

func (rp *DynamicReverseProxy) proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("proxy error", "error", err)
	http.Error(w, "upstream error", http.StatusBadGateway)
}

func (rp *DynamicReverseProxy) buildReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		networktools.PrepareProxyHeaders(req, target)
	}

	proxy.ErrorHandler = rp.proxyErrorHandler
	return proxy
}

func (rp *DynamicReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), rp.timeout)
	defer cancel()
	cleanHost := strings.Split(r.Host, ":")[0]
	target, err := rp.resolveServiceTarget(ctx, cleanHost)
	if err != nil {
		slog.Warn("target resolution failed", "host", r.Host, "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	proxy := rp.buildReverseProxy(target)
	handler := http.TimeoutHandler(proxy, rp.timeout, "proxy timeout")

	protoHeader := r.Header.Get("X-Forwarded-Proto")
	protoParts := strings.Split(protoHeader, ",")
	proto := ""
	if len(protoParts) > 0 {
		proto = strings.ToLower(strings.TrimSpace(protoParts[0]))
	}
	isRequestSecure := r.TLS != nil || proto == "https"

	if isRequestSecure || !rp.useHttps {
		handler.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	// TODO: check if the host has HTTPS certificates

	redirectHost := cleanHost
	if rp.httpsPort != "" && rp.httpsPort != "443" {
		redirectHost = fmt.Sprintf("%s:%s", cleanHost, rp.httpsPort)
	}
	redirectUrl := fmt.Sprintf("https://%s%s", redirectHost, r.URL.RequestURI())
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}
