package udp

import (
	"fmt"
	"github.com/FDUTCH/obsidian/internal/pool"
	"github.com/FDUTCH/obsidian/internal/util"
	"github.com/FDUTCH/obsidian/proxy/balance"
	"net"
	"sync"
	"time"
)

// Proxy - UDP implementation of proxy.Proxy
type Proxy struct {
	connections *sync.Map
	pool        *pool.Pool[[]byte]
	network     string
}

func NewProxy(buffSize int, network string) (*Proxy, error) {
	if buffSize < 1 {
		return nil, fmt.Errorf("bufferSize can not be of length less than 1")
	}
	return &Proxy{pool: util.NewBufferPool(buffSize), connections: new(sync.Map), network: network}, nil
}

func (p *Proxy) Listen(remoteAddress, localAddress string) error {

	conn, err := net.ListenPacket(p.network, localAddress)
	if err != nil {
		return err
	}
	remote, err := net.ResolveUDPAddr(p.network, remoteAddress)
	if err != nil {
		return err
	}

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
			c, err := dial(p.network, remote.String())
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

func (p *Proxy) Balance(balancer balance.LoadBalancer, localAddress string) error {
	conn, err := net.ListenPacket("udp", localAddress)
	if err != nil {
		return err
	}

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
		})

		session := w.(*Session)

		if !loaded {
			dial := (&net.Dialer{}).Dial
			c, err := dial(p.network, balancer.Address())
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
				if time.Since(session.LastActivity()) > time.Minute {
					defer session.Close()
				}
				return true
			})
		}
	}
}
