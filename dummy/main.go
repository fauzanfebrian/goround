package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	// ports := []string{"8081", "8082", "8083"}
	ports := []string{"8081", "8083"}
	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("Received request on %s\n", p)
				fmt.Fprintf(w, "Hello from Server %s\n", p)
			})

			log.Printf("Server listening on port %s", p)
			if err := http.ListenAndServe(":"+p, mux); err != nil {
				log.Fatal(err)
			}
		}(port)
	}

	wg.Wait()
}
