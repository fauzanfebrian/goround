package pool

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ServerPool struct {
	Port  int
	Url   *url.URL
	alive bool
	mu    sync.RWMutex
}

// Server checking using HTTP Method / Layer 7
func (serverPool *ServerPool) CheckServerHttp() {
	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

	resp, err := http.Get(serverPool.Url.String())
	if err != nil {
		fmt.Printf("Check server %d error: %s\n", serverPool.Port, err)
		serverPool.alive = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Check server %d error: Got statusCode %d from server\n", serverPool.Port, resp.StatusCode)
		return
	}

	fmt.Printf("Check server %d succeed\n", serverPool.Port)
	serverPool.alive = true
}

// Server checking using TCP Connection / Layer 4 (From instruction)
func (serverPool *ServerPool) CheckServer() {
	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

	address := fmt.Sprintf("localhost:%d", serverPool.Port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Printf("Check server %d error: %s\n", serverPool.Port, err)
		serverPool.alive = false
		return
	}
	defer conn.Close()

	fmt.Printf("Check server %d succeed\n", serverPool.Port)
	serverPool.alive = true
}

func (serverPool *ServerPool) IsAlive() bool {
	serverPool.mu.RLock()
	defer serverPool.mu.Unlock()
	return serverPool.alive
}
