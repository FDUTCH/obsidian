package udp

import (
	"github.com/FDUTCH/obsidian/proxy/balance"
	"github.com/FDUTCH/obsidian/proxy/packet"
)

// Proxy - UDP implementation of proxy.Proxy
type Proxy struct {
	packetProxy *packet.Proxy
}

func NewProxy(buffSize int, network string) (*Proxy, error) {
	p, err := packet.NewProxy(buffSize, network)
	if err != nil {
		return nil, err
	}
	return &Proxy{packetProxy: p}, nil
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {
	return p.packetProxy.Listen(remoteAddress, localAddress)
}

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {
	return p.packetProxy.Balance(balancer, localAddress)
}
