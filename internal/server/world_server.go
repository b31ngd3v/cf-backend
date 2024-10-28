package server

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/inconshreveable/muxado"
)

type worldServer struct {
	addr   string
	sess   muxado.Session
	quitCh chan struct{}
}

func NewWorldServer(port int, conn net.Conn) *worldServer {
	sess := muxado.Server(conn, nil)
	return &worldServer{
		addr:   fmt.Sprintf(":%d", port),
		sess:   sess,
		quitCh: make(chan struct{}),
	}
}

func (s *worldServer) Run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return
	}
	defer listener.Close()
	defer close(s.quitCh)

	go func() {
		for {
			select {
			case <-s.quitCh:
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					continue
				}
				go s.handleConn(conn)
			}
		}
	}()

	s.sess.Wait()
}

func (s *worldServer) handleConn(conn net.Conn) {
	defer conn.Close()
	stream, err := s.sess.OpenStream()
	if err != nil {
		return
	}
	defer stream.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 2)

	go s.forwardPackets(ctx, stream, conn, errCh)
	go s.forwardPackets(ctx, conn, stream, errCh)

	select {
	case <-s.quitCh:
	case <-errCh:
		cancel()
	}
}

func (s *worldServer) forwardPackets(ctx context.Context, dst io.Writer, src io.Reader, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			d := make([]byte, bufferSize)
			bc, err := src.Read(d)
			if err == io.EOF {
				errCh <- io.EOF
				return
			} else if err != nil {
				continue
			}
			_, err = dst.Write(d[:bc])
			if err == io.EOF {
				errCh <- io.EOF
				return
			} else if err != nil {
				continue
			}
		}
	}
}
