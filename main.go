package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	b "github.com/fauzanfebrian/goround/backend"
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

	port := 8000
	addr := fmt.Sprintf(":%d", port)
	backend := b.NewBackend(serverPools)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reverseProxy := backend.GetReverseProxy()
		if reverseProxy != nil {
			reverseProxy.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "There's no available server right now, try again later\n")
	})

	log.Printf("Server listening on port %d", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
