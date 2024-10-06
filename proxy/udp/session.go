package udp

import (
	"github.com/FDUTCH/obsidian/internal/pool"
	"io"
	"net"
	"sync"
	"time"
)

// Session implements active udp session
type Session struct {
	conn         net.PacketConn
	addr         net.Addr
	pool         *pool.Pool[[]byte]
	buff         []byte
	mp           *sync.Map
	writer       io.WriteCloser
	mu           *sync.Mutex
	lastActivity time.Time
}

func (s *Session) Close() error {
	s.pool.Put(s.buff)
	s.mp.Delete(s.addr.String())
	_ = s.writer.Close()
	return nil
}

func (s *Session) Write(p []byte) (n int, err error) {
	return s.writer.Write(p)
}

func (s *Session) handle(conn net.Conn) {

	s.writer = conn

	s.buff = s.pool.Get()

	go func() {
		defer s.Close()

		for {
			n, err := conn.Read(s.buff)
			if err != nil {
				return
			}
			s.active()
			_, err = s.conn.WriteTo(
				s.buff[:n],
				s.addr)
			if err != nil {
				return
			}
		}
	}()
}

func (s *Session) active() {
	s.mu.Lock()
	s.lastActivity = time.Now()
	s.mu.Unlock()
}

func (s *Session) LastActivity() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastActivity
}
