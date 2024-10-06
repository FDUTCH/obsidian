package http_proxy

import (
	"github.com/FDUTCH/obsidian/proxy/balance"
	"net/http"
)

// BalanceHandler http.Handler implementation for balancing http load
type BalanceHandler struct {
	balancer balance.LoadBalancer
	c        *cache
}

func NewBalanceHandler(balancer balance.LoadBalancer) *BalanceHandler {
	return &BalanceHandler{balancer: balancer, c: new(cache)}
}

func (b *BalanceHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	proxy, err := b.c.getProxy(b.balancer.Address())
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	proxy.ServeHTTP(writer, request)
}
