package http_proxy

import (
	"github.com/FDUTCH/obsidian/proxy/balance"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy - http implementation of proxy.Proxy
type Proxy struct {
	keyFile, certFile string
}

func NewProxy() *Proxy {
	return &Proxy{}
}

func NewSecureProxy(keyFile, certFile string) *Proxy {
	return &Proxy{
		keyFile:  keyFile,
		certFile: certFile,
	}
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {
	uri, err := newUrl(remoteAddress)
	if err != nil {
		return err
	}
	return p.serve(localAddress, httputil.NewSingleHostReverseProxy(uri))
}

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {
	handler := NewBalanceHandler(balancer)

	return p.serve(localAddress, handler)
}

func (p *Proxy) Split(splitter RouteSplitter, localAddress string) error {
	handler := &RouteSplitterHandler{splitter: splitter}

	return p.serve(localAddress, handler)
}

func (p *Proxy) serve(local string, handler http.Handler) error {
	if p.keyFile == "" || p.certFile == "" {
		return http.ListenAndServe(local, handler)
	}
	return http.ListenAndServeTLS(local, p.certFile, p.keyFile, handler)
}

func newUrl(address string) (*url.URL, error) {
	uri, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	if uri.Scheme == "" {
		uri, err = url.Parse("https://" + address + "/")
	}
	return uri, err
}
