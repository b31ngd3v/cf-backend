package server

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestPerformPawshake(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	deadline := time.Now().Add(100 * time.Millisecond)
	clientConn.SetDeadline(deadline)

	var wg sync.WaitGroup
	wg.Add(1)

	port := 12345
	go func() {
		defer wg.Done()
		clientConn.Write([]byte(pawshakeRequest))
		buf := make([]byte, bufferSize)
		n, err := clientConn.Read(buf)
		if err != nil {
			t.Errorf("Client read error: %v", err)
		}
		response := string(buf[:n])
		expectedResponse := fmt.Sprintf("%s%d\n", pawshakeRequestSucceeded, port)
		if response != expectedResponse {
			t.Errorf("Expected response %q, got %q", expectedResponse, response)
		}
	}()

	cs := &clientServer{}
	success := cs.performPawshake(serverConn, port)
	if !success {
		t.Errorf("performPawshake failed when it should have succeeded")
	}

	wg.Wait()
}

func TestPerformPawshakeFail(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	deadline := time.Now().Add(100 * time.Millisecond)
	clientConn.SetDeadline(deadline)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		clientConn.Write([]byte("invalid_request"))
		buf := make([]byte, bufferSize)
		n, err := clientConn.Read(buf)
		if err != nil {
			t.Errorf("Client read error: %v", err)
		}
		response := string(buf[:n])
		if response != pawshakeRequestFailed {
			t.Errorf("Expected response %q, got %q", pawshakeRequestFailed, response)
		}
	}()

	port := 12345
	cs := &clientServer{}
	success := cs.performPawshake(serverConn, port)
	if success {
		t.Errorf("performPawshake succeeded when it should have failed")
	}

	wg.Wait()
}
