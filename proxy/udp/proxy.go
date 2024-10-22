package udp

import (
	"fmt"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"github.com/FDUTCH/obsidian/proxy/packet"
)

// Proxy - UDP implementation of proxy.Proxy
type Proxy struct {
	network  string
	buffSize int
}

func NewProxy(buffSize int, network string) (*Proxy, error) {
	if buffSize < 1 {
		return nil, fmt.Errorf("bufferSize can not be of length less than 1")
	}
	return &Proxy{network: network, buffSize: buffSize}, nil
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {
	pr, err := packet.NewProxy(p.buffSize, p.network)
	if err != nil {
		return err
	}
	return pr.Listen(remoteAddress, localAddress)
}

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {
	pr, err := packet.NewProxy(p.buffSize, p.network)
	if err != nil {
		return err
	}
	return pr.Balance(balancer, localAddress)
}
