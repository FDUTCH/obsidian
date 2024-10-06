package proxy

import (
	"fmt"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"github.com/FDUTCH/obsidian/proxy/http_proxy"
	"github.com/FDUTCH/obsidian/proxy/tcp"
	"github.com/FDUTCH/obsidian/proxy/udp"
	"net"
)

// Proxy ...
type Proxy interface {
	// Listen redirects incoming traffic
	Listen(remoteAddress, localAddress string) error
	// Balance balances incoming traffic using balance.LoadBalancer
	Balance(balancer balance.LoadBalancer, localAddress string) error
}

// HttpProxy ...
type HttpProxy interface {
	Proxy
	// Split splits incoming traffic using http_proxy.RouteSplitter
	Split(splitter http_proxy.RouteSplitter, localAddress string) error
}

// Config ...
type Config interface {
	Run()
}

// New makes a new Proxy
func New(network string, bufferSize int, params ...string) (Proxy, error) {
	switch network {
	case "tcp", "tcp4":
		return tcp.NewProxy(bufferSize, network)
	case "udp", "udp4", "udp6":
		return udp.NewProxy(bufferSize, network)
	case "http":
		return http_proxy.NewProxy(), nil
	case "https":
		if len(params) != 2 {
			return nil, fmt.Errorf("invalid param count")
		}
		return http_proxy.NewSecureProxy(params[0], params[1]), nil
	}
	return nil, net.UnknownNetworkError(network)
}
