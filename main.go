package main

import (
	"fmt"
	"time"

	"github.com/fauzanfebrian/goround/pool"
)

func main() {
	serverPools := pool.CreateServerPools(8081, 8082, 8083)

	for !serverPools[0].IsAlive() && !serverPools[1].IsAlive() && !serverPools[2].IsAlive() {
		fmt.Println("Waiting for a pool up and runing")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Load balancer is ready to accept a connection")
}
