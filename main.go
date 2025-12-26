package main

import (
	"fmt"
	"time"

	"github.com/fauzanfebrian/goround/pool"
)

func main() {
	serverPools := pool.CreateServerPools(8081, 8082, 8083)

	for !serverPools[0].Alive && !serverPools[1].Alive && !serverPools[2].Alive {
		fmt.Println("Waiting for a pool up and runing")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("Load balancer is ready to accept a connection")
}
