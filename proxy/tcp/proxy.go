package tcp

import (
	"fmt"
	"github.com/FDUTCH/obsidian/internal/pool"
	"github.com/FDUTCH/obsidian/internal/util"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"io"
	"net"
)

// Proxy - TCP implementation of proxy.Proxy
type Proxy struct {
	pool    *pool.Pool[[]byte]
	network string
}

func NewProxy(buffSize int, network string) (*Proxy, error) {
	if buffSize < 1 {
		return nil, fmt.Errorf("bufferSize can not be of length less than 1")
	}
	return &Proxy{pool: util.NewBufferPool(buffSize), network: network}, nil
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {
	remote, err := net.ResolveTCPAddr(p.network, remoteAddress)
	if err != nil {
		return err
	}

	listener, err := net.Listen(p.network, localAddress)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		handleConn(conn, remote.String(), p.pool, p.network)
	}

}

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {

	listener, err := net.Listen(p.network, localAddress)

	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go handleConn(conn, balancer.Address(), p.pool, p.network)
	}
}

func handleConn(conn net.Conn, addr string, provider *pool.Pool[[]byte], network string) {
	serverConn, err := net.Dial(network, addr)
	if err != nil {
		_ = conn.Close()
		return
	}

	go func() {
		defer serverConn.Close()
		_, _ = io.CopyBuffer(serverConn, conn, provider.Get())
	}()

	go func() {
		defer conn.Close()
		_, _ = io.CopyBuffer(conn, serverConn, provider.Get())
	}()

}
