package proxy

import (
	"github.com/FDUTCH/obsidian/internal/util"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"github.com/FDUTCH/obsidian/proxy/http_proxy"
)

// Options represents the options that control how the proxy should be set up.
type Options struct {
	// Port - port that will be listened to for new connections
	Port int `json:"port"`

	// Network - network that will be used for transfer data
	Network string `json:"network"`

	// BufferSize - size of the buffer that will be used for transferring data
	BufferSize int `json:"buffer_size,omitempty"`

	// Router - maps http route to address of the remote server (can be used only for http/https)
	Router map[string]string `json:"router,omitempty"`

	// RemoteAddress - the address to which traffic will be redirected
	RemoteAddress string `json:"remote_address,omitempty"`

	// Servers - list off servers between which will be balanced incoming traffic
	Servers []string `json:"servers,omitempty"`

	// Key - used for https configuration
	Key string `json:"key,omitempty"`

	// Cert - used for https configuration
	Cert string `json:"cert,omitempty"`
}

// Run creates and runs proxy
func (c Options) Run() error {
	p, err := New(c.Network, c.BufferSize, c.Key, c.Cert)
	if err != nil {
		return err
	}
	switch {

	case len(c.Servers) > 0:
		// LoadBalancer mode
		return p.Balance(balance.NewSimpleLoadBalancer(c.Servers...), util.Addr(c.Port))
	case len(c.Router) > 0:
		// RouteSplitter mode
		return p.(*http_proxy.Proxy).Split(http_proxy.NewRouteSplitter(c.Router), util.Addr(c.Port))
	default:
		// Default mode
		return p.Listen(c.RemoteAddress, util.Addr(c.Port))
	}
}
