package pool

import (
	"fmt"
	"net/http"
	"sync"
)

type ServerPool struct {
	Port  int
	Url   string
	alive bool
	mu    sync.RWMutex
}

func (serverPool *ServerPool) CheckServer() {
	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()

	resp, err := http.Get(serverPool.Url)
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

func (serverPool *ServerPool) IsAlive() bool {
	serverPool.mu.Lock()
	defer serverPool.mu.Unlock()
	return serverPool.alive
}
