package pool

import (
	"fmt"
	"net/url"
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
		rawUrl := fmt.Sprintf("http://localhost:%d", port)

		url, err := url.Parse(rawUrl)
		if err != nil {
			fmt.Printf("URL parse error for %d: %s", port, err)
			continue
		}

		serverPool := ServerPool{
			Port: port,
			Url:  url,
		}

		go checkServer(&serverPool)

		serverPools = append(serverPools, &serverPool)
	}

	return serverPools
}
