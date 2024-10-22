package packet

import (
	"fmt"
	"github.com/FDUTCH/obsidian/internal/pool"
	"github.com/FDUTCH/obsidian/internal/util"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"net"
	"sync"
	"time"
)

// Proxy - PacketConn implementation of proxy.Proxy
type Proxy struct {
	connections *sync.Map
	pool        *pool.Pool[[]byte]
	network     string
}

func NewProxy(buffSize int, network string) (*Proxy, error) {
	if buffSize < 1 {
		return nil, fmt.Errorf("bufferSize can not be less than 1")
	}
	return &Proxy{pool: util.NewBufferPool(buffSize), connections: new(sync.Map), network: network}, nil
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {
	conn, err := net.ListenPacket(p.network, localAddress)
	if err != nil {
		return err
	}
	addr, err := net.ResolveIPAddr("", remoteAddress)

	remote := addr.String()

	if err != nil {
		return err
	}

	return p.listen(func() string { return remote }, conn)
}

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {
	conn, err := net.ListenPacket(p.network, localAddress)
	if err != nil {
		return err
	}

	return p.listen(balancer.Address, conn)
}

func (p *Proxy) listen(remote func() (remoteAddress string), conn net.PacketConn) error {
	local := conn.LocalAddr()

	buff := p.pool.Get()

	closed := make(chan struct{})

	defer close(closed)

	go p.gc(closed)

	for {
		n, addr, err := conn.ReadFrom(buff)
		if err != nil {
			return err
		}
		w, loaded := p.connections.LoadOrStore(addr.String(), &Session{
			addr: addr,
			pool: p.pool,
			mp:   p.connections,
			conn: conn,
			mu:   &sync.Mutex{},
		})

		session := w.(*Session)

		if !loaded {
			dial := (&net.Dialer{}).Dial
			c, err := dial(local.Network(), remote())
			if err != nil {
				return err
			}
			session.handle(c)
		}

		_, err = session.Write(buff[:n])

		if err != nil {
			_ = session.Close()
		}

	}
}

// gc clearing all outdated connections
func (p *Proxy) gc(closed chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-closed:
			return
		case <-ticker.C:
			p.connections.Range(func(key, value any) bool {
				session := value.(*Session)
				if time.Since(session.LastActivity()) > time.Second*30 {
					defer session.Close()
				}
				return true
			})
		}
	}
}
