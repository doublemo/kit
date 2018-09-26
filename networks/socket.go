package networks

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Socket 实现TCP服务
type Socket struct {
	name     string           //
	Config   *SocketConfig    // 配置
	status   int              //
	done     chan error       //
	quit     chan struct{}    // 退出信号
	wg       sync.WaitGroup   //
	mux      sync.Mutex       //
	listener *net.TCPListener //
}

// Serve 启动服务
func (s *Socket) Serve() error {
	defer func() {
		close(s.done)
	}()

	s.mux.Lock()
	if s.status != SocketStateStoped {
		return ErrorServerNonStoped
	}

	s.status = SocketStateRunning
	s.mux.Unlock()

	go s.serve()
	err := <-s.done
	close(s.quit)
	s.listener.SetDeadline(time.Now())
	s.listener.Close()
	s.wg.Wait()

	s.mux.Lock()
	s.status = SocketStateStoped
	s.mux.Unlock()
	return err
}

func (s *Socket) serve() {
	if err := s.listen(s.Config.Addr); err != nil {
		s.done <- err
		return
	}

	if s.Config.CallBack == nil {
		s.done <- fmt.Errorf("CallBack is undefined")
		return
	}

	log.Println("TCP listening on:", s.listener.Addr())
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			return
		}

		conn.SetReadBuffer(s.Config.ReadBufferSize)
		conn.SetWriteBuffer(s.Config.WriteBufferSize)
		go s.client(conn)
	}
}

func (s *Socket) listen(addr string) error {
	resolveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	lis, err := net.ListenTCP("tcp", resolveAddr)
	if err != nil {
		return err
	}

	s.listener = lis
	return nil
}

func (s *Socket) client(conn net.Conn) {
	defer func() {
		s.wg.Done()
		conn.Close()
	}()

	s.wg.Add(1)
	s.Config.CallBack(conn, s.quit)
}

// Shutdown 关闭socket 服务
func (s *Socket) Shutdown() {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.status != SocketStateRunning {
		return
	}
	s.done <- nil
}

// NewSocket 创建socket服务
func NewSocket(c *SocketConfig) *Socket {
	return &Socket{
		done:   make(chan error),
		quit:   make(chan struct{}),
		Config: c,
	}
}
