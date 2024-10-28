package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

const bufferSize = 1024

var (
	c                        = GetConfig()
	pawshakeTimeout          = 5 * time.Second
	versionString            = fmt.Sprintf("cat_forwarding_v%s", c.Version)
	maxConnLimitExceeded     = fmt.Sprintf("%s\nmax_conn_limit_exceeded", versionString)
	pawshakeRequest          = fmt.Sprintf("%s\npawshake_request\n", versionString)
	pawshakeRequestFailed    = fmt.Sprintf("%s\nno_pawshake_w_bad_cat\n", versionString)
	pawshakeRequestSucceeded = fmt.Sprintf("%s\npawshake_successful\nport_", versionString)
)

func Run() {
	s := newClientServer(c.ClientPort)
	s.run()
}

type clientServer struct {
	addr string
}

func newClientServer(port int) *clientServer {
	return &clientServer{
		addr: fmt.Sprintf(":%d", port),
	}
}

func (s *clientServer) run() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalln("client port already in use!")
	}
	log.Printf("client server is listening on %s\n", s.addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *clientServer) handleConn(conn net.Conn) {
	defer conn.Close()
	port, err := AssignPort(&conn)
	if err != nil {
		conn.Write([]byte(maxConnLimitExceeded))
		return
	}
	defer FreePort(port)
	if !s.performPawshake(conn, port) {
		return
	}
	ws := NewWorldServer(port, conn)
	ws.Run()
}

func (s *clientServer) performPawshake(conn net.Conn, port int) bool {
	conn.SetDeadline(time.Now().Add(pawshakeTimeout))
	defer conn.SetDeadline(time.Time{})
	buffer := make([]byte, bufferSize)
	bc, err := conn.Read(buffer)
	if err != nil {
		return false
	}
	if pawshakeRequest != string(buffer[:bc]) {
		conn.Write([]byte(pawshakeRequestFailed))
		return false
	}
	res := fmt.Sprintf("%s%d\n", pawshakeRequestSucceeded, port)
	conn.Write([]byte(res))
	return true
}
