package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/inconshreveable/muxado"
)

func TestForwardPackets(t *testing.T) {
	errCh := make(chan error, 1)

	testStr := "hello world"
	src := bytes.NewBufferString(testStr)
	dst := &bytes.Buffer{}

	ws := &worldServer{}
	ws.forwardPackets(context.TODO(), dst, src, errCh)

	if !bytes.Equal([]byte(testStr), dst.Bytes()) {
		t.Errorf("Expected data %q, got %q", testStr, dst)
	}
}

func TestWorldServerRun(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	port := 12345
	ws := NewWorldServer(port, serverConn)
	defer ws.sess.Close()

	go ws.Run()

	time.Sleep(100 * time.Millisecond)

	timeout := 100 * time.Millisecond
	conn, err := net.DialTimeout("tcp", ws.addr, timeout)
	if err != nil {
		t.Fatalf("Failed to connect to worldServer: %v", err)
	}
	defer conn.Close()

	deadline := time.Now().Add(timeout)
	conn.SetDeadline(deadline)

	message := []byte("Test message")
	_, err = conn.Write(message)
	if err != nil {
		t.Fatalf("Failed to write to worldServer: %v", err)
	}

	clientConn.SetDeadline(deadline)

	sess := muxado.Client(clientConn, nil)
	defer sess.Close()
	stream, err := sess.Accept()
	if err != nil {
		t.Fatalf("Failed to read from clientConn: %v", err)
	}

	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		t.Fatalf("Failed to read from clientConn: %v", err)
	}

	if !bytes.Equal(buf, message) {
		t.Errorf("Expected message %q, got %q", message, buf)
	}
}
