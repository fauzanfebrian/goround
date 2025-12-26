package pool

import (
	"fmt"
	"net/http"
)

type ServerPool struct {
	Port  int
	Alive bool
}

func (pool *ServerPool) CheckServer() {
	url := fmt.Sprintf("http://localhost:%d", pool.Port)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Check server %d error: %s\n", pool.Port, err)
		pool.Alive = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Check server %d error: Got statusCode %d from server\n", pool.Port, resp.StatusCode)
		return
	}

	fmt.Printf("Check server %d succeed\n", pool.Port)
	pool.Alive = true
}
