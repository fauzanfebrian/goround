package pool

import (
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
		serverPool := ServerPool{
			Port: port,
		}

		go checkServer(&serverPool)

		serverPools = append(serverPools, &serverPool)
	}

	return serverPools
}
