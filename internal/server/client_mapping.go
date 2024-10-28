package server

import (
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"
)

var (
	mu        sync.RWMutex
	clientMap = make(map[int]*net.Conn)
	connPerIP = make(map[string]int)
	randSeed  = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func AssignPort(conn *net.Conn) (int, error) {
	mu.Lock()
	defer mu.Unlock()
	ip := getIPFromConn(conn)
	if connPerIP[ip] >= c.ConnLimit {
		return 0, errors.New("max_conn_limit_exceeded")
	}
	for range 10_000 {
		port := getRandomPort()
		if port == c.ClientPort {
			continue
		}
		_, assigned := clientMap[port]
		if !assigned {
			clientMap[port] = conn
			connPerIP[ip]++
			return port, nil
		}
	}
	return 0, errors.New("all_ports_used")
}

func GetPort(port int) (*net.Conn, bool) {
	mu.RLock()
	defer mu.RUnlock()
	val, ok := clientMap[port]
	return val, ok
}

func FreePort(port int) {
	mu.Lock()
	defer mu.Unlock()
	conn := clientMap[port]
	ip := getIPFromConn(conn)
	connPerIP[ip]--
	if connPerIP[ip] <= 0 {
		delete(connPerIP, ip)
	}
	delete(clientMap, port)
}

func getIPFromConn(conn *net.Conn) string {
	remoteAddr := (*conn).RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(remoteAddr)
	return ip
}

func getRandomPort() int {
	// 64512 = 65535 (max) - 1024 (min) + 1
	return randSeed.Intn(64512) + 1024
}
