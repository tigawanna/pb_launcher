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
	http01 "pb_launcher/internal/certificates/http_01"
	launcherdomain "pb_launcher/internal/launcher/domain"
	proxydomain "pb_launcher/internal/proxy/domain"
	"pb_launcher/utils/networktools"
	"strconv"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/apis"
)

type DynamicReverseProxy struct {
	discovery           *proxydomain.ServiceDiscovery
	domainDiscovery     *proxydomain.DomainServiceDiscovery
	installTokenUsecase *launcherdomain.CleanServiceInstallTokenUsecase
	apiDomain           string
	apiAddress          string
	useHttps            bool
	skipHttpsRedirect   bool
	httpsPort           string
	timeout             time.Duration
	http01Store         *http01.Http01ChallengeAddressPublisher
}

var _ http.Handler = (*DynamicReverseProxy)(nil)

func NewDynamicReverseProxy(
	discovery *proxydomain.ServiceDiscovery,
	domainDiscovery *proxydomain.DomainServiceDiscovery,
	installTokenUsecase *launcherdomain.CleanServiceInstallTokenUsecase,
	http01Store *http01.Http01ChallengeAddressPublisher,
	cfg configs.Config,
	pbConf *apis.ServeConfig,
) *DynamicReverseProxy {
	return &DynamicReverseProxy{
		discovery:           discovery,
		domainDiscovery:     domainDiscovery,
		installTokenUsecase: installTokenUsecase,
		http01Store:         http01Store,
		apiDomain:           cfg.GetDomain(),
		useHttps:            cfg.IsHttpsEnabled(),
		skipHttpsRedirect:   cfg.IsHttpsRedirectDisabled(),
		httpsPort:           cfg.GetHttpsPort(),
		apiAddress:          pbConf.HttpAddr,
		timeout:             15 * time.Second,
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

const superusersEndpoint = "/api/collections/_superusers/records"

func (rp *DynamicReverseProxy) proxyModifyResponse(r *http.Response) error {
	if r.Request.Method == http.MethodPost &&
		strings.HasPrefix(r.Request.URL.Path, superusersEndpoint) &&
		r.StatusCode == 200 {
		authorization := r.Request.Header.Get("Authorization")
		rp.installTokenUsecase.CleanInstallToken(r.Request.Context(), authorization)
	}
	return nil
}

func (rp *DynamicReverseProxy) buildReverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		networktools.PrepareProxyHeaders(req, target)
	}
	proxy.ModifyResponse = rp.proxyModifyResponse
	proxy.ErrorHandler = rp.proxyErrorHandler
	return proxy
}

const AcmeChallengePath = "/.well-known/acme-challenge/"

func (rp *DynamicReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), rp.timeout)
	defer cancel()
	cleanHost := strings.Split(r.Host, ":")[0]

	var proxy *httputil.ReverseProxy

	isAcmeChallenge := strings.HasPrefix(r.URL.Path, AcmeChallengePath)
	if isAcmeChallenge {
		targetURL, err := rp.http01Store.ResolveAddress()
		if err != nil {
			http.Error(w, "not found", http.StatusInternalServerError)
			return
		}
		proxy = httputil.NewSingleHostReverseProxy(targetURL) // For some reason, Let's Encrypt doesn't seem to work well with my buildReverseProxy
	} else {
		targetURL, err := rp.resolveServiceTarget(ctx, cleanHost)
		if err != nil || targetURL == nil {
			slog.Warn("target resolution failed", "host", r.Host, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		proxy = rp.buildReverseProxy(targetURL)
	}

	handler := http.TimeoutHandler(proxy, rp.timeout, "proxy timeout")

	if networktools.IsRequestSecure(r) || !rp.useHttps || rp.skipHttpsRedirect || isAcmeChallenge {
		handler.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	redirectUrl := networktools.BuildHostURL("https", cleanHost, rp.httpsPort, r.URL.RequestURI())
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}
