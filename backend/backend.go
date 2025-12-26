package backend

import (
	"fmt"
	"net/http/httputil"
	"sync/atomic"

	"github.com/fauzanfebrian/goround/pool"
)

type Backend struct {
	poolsLen    int
	counter     int64
	proxies     []*httputil.ReverseProxy
	serverPools []*pool.ServerPool
}

// get index from live server
func (backend *Backend) getLiveIndex() int64 {
	counter := atomic.AddInt64(&backend.counter, 1)

	index := (counter - 1) % int64(backend.poolsLen)
	if backend.serverPools[index].IsAlive() {
		return index
	}

	for range backend.poolsLen {
		counter = atomic.AddInt64(&backend.counter, 1)
		index = (counter - 1) % int64(backend.poolsLen)

		fmt.Printf("Index: %d and Counter: %d\n", index, counter)

		if backend.serverPools[index].IsAlive() {
			return index
		}
	}

	return -1
}

// When this function got triggered it's mean the request occured,
// therefore the counter is incremented.
func (backend *Backend) GetReverseProxy() *httputil.ReverseProxy {
	index := backend.getLiveIndex()

	if index < 0 {
		return nil
	}

	return backend.proxies[index]
}
