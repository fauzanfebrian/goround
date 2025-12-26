package pool

import (
	"fmt"
	"time"
)

func checkServer(serverPool *ServerPool) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		serverPool.CheckServer()
	}
}

func CreateServerPools(ports ...int) []*ServerPool {
	serverPools := []*ServerPool{}

	for _, port := range ports {
		url := fmt.Sprintf("http://localhost:%d", port)
		serverPool := ServerPool{
			Port: port,
			Url:  url,
		}

		go checkServer(&serverPool)

		serverPools = append(serverPools, &serverPool)
	}

	return serverPools
}
