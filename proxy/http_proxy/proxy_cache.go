package http_proxy

import (
	"net/http/httputil"
	"sync"
)

// cache caches [address => proxy] pair
type cache struct {
	proxies sync.Map
}

// getProxy returns proxy by address
func (c *cache) getProxy(addr string) (*httputil.ReverseProxy, error) {
	p, ok := c.proxies.Load(addr)
	if !ok {
		uri, err := newUrl(addr)
		if err != nil {
			return nil, err
		}
		p = httputil.NewSingleHostReverseProxy(uri)
		c.proxies.Store(addr, p)
	}
	return p.(*httputil.ReverseProxy), nil
}
