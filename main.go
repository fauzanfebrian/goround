package main

import (
	"fmt"
	"time"

	"github.com/fauzanfebrian/goround/pool"
)

func main() {
	serverPools := pool.CreateServerPools(8081, 8082, 8083)

	fmt.Println("Waiting for a pool up and runing")
	for true {
		isAPoolAlive := false

		for _, serverPool := range serverPools {
			if serverPool.IsAlive() {
				isAPoolAlive = true
				break
			}
		}

		if isAPoolAlive {
			break
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("Load balancer is ready to accept a connection")
}
