package server

import (
	"io"
	"net"
	"testing"
	"time"
)

type MockAddr struct {
	network string
	address string
}

func (m *MockAddr) Network() string { return m.network }
func (m *MockAddr) String() string  { return m.address }

type MockConn struct {
	remoteAddr net.Addr
}

func (m *MockConn) Read(b []byte) (n int, err error)   { return 0, io.EOF }
func (m *MockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return m.remoteAddr }
func (m *MockConn) RemoteAddr() net.Addr               { return m.remoteAddr }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestAssignPort(t *testing.T) {
	remoteAddr := &MockAddr{
		network: "tcp",
		address: "192.168.1.1:12345",
	}

	conn := &MockConn{
		remoteAddr: remoteAddr,
	}
	connPtr := net.Conn(conn)

	port, err := AssignPort(&connPtr)
	if err != nil {
		t.Errorf("AssignPort failed: %v", err)
	}

	if port == c.ClientPort {
		t.Errorf("Assigned port should not be equal to ClientPort")
	}

	conn2 := &MockConn{
		remoteAddr: remoteAddr,
	}
	connPtr2 := net.Conn(conn2)
	port2, err := AssignPort(&connPtr2)
	if err != nil {
		t.Errorf("AssignPort failed: %v", err)
	}

	if port2 == c.ClientPort || port2 == port {
		t.Errorf("Assigned port should not be equal to ClientPort or previous port")
	}

	conn3 := &MockConn{
		remoteAddr: remoteAddr,
	}
	connPtr3 := net.Conn(conn3)
	_, err = AssignPort(&connPtr3)
	if err == nil {
		t.Errorf("AssignPort should fail due to ConnLimit")
	}

	FreePort(port)

	conn4 := &MockConn{
		remoteAddr: remoteAddr,
	}
	connPtr4 := net.Conn(conn4)
	port4, err := AssignPort(&connPtr4)
	if err != nil {
		t.Errorf("AssignPort failed: %v", err)
	}

	if port4 == c.ClientPort || port4 == port2 {
		t.Errorf("Assigned port should not be equal to ClientPort or previous ports")
	}
}
