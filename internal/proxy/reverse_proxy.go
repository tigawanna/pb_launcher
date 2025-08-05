package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"pb_launcher/configs"
	http01 "pb_launcher/internal/certificates/http_01"
	"pb_launcher/utils/networktools"
	"strings"
	"time"
)

type DynamicReverseProxy struct {
	proxyResolver     *DynamicReverseProxyDiscovery
	useHttps          bool
	skipHttpsRedirect bool
	httpsPort         string
	timeout           time.Duration
	http01Store       *http01.Http01ChallengeAddressPublisher
}

var _ http.Handler = (*DynamicReverseProxy)(nil)

func NewDynamicReverseProxy(
	proxyResolver *DynamicReverseProxyDiscovery,
	http01Store *http01.Http01ChallengeAddressPublisher,
	cfg configs.Config,
) *DynamicReverseProxy {
	return &DynamicReverseProxy{
		proxyResolver:     proxyResolver,
		http01Store:       http01Store,
		useHttps:          cfg.IsHttpsEnabled(),
		skipHttpsRedirect: cfg.IsHttpsRedirectDisabled(),
		httpsPort:         cfg.GetHttpsPort(),
		timeout:           15 * time.Second,
	}
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
		var err error
		proxy, err = rp.proxyResolver.ResolveTarget(ctx, cleanHost)
		if err != nil || proxy == nil {
			slog.Warn("target resolution failed", "host", r.Host, "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	handler := http.TimeoutHandler(proxy, rp.timeout, "proxy timeout")

	if networktools.IsRequestSecure(r) || !rp.useHttps || rp.skipHttpsRedirect || isAcmeChallenge {
		handler.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	redirectUrl := networktools.BuildHostURL("https", cleanHost, rp.httpsPort, r.URL.RequestURI())
	http.Redirect(w, r, redirectUrl, http.StatusPermanentRedirect)
}
