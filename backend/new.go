package backend

import (
	"net/http/httputil"

	"github.com/fauzanfebrian/goround/pool"
)

func NewBackend(serverPools []*pool.ServerPool) *Backend {
	poolsLen := len(serverPools)

	var proxies []*httputil.ReverseProxy
	for _, serverPool := range serverPools {
		reverseProxy := httputil.NewSingleHostReverseProxy(serverPool.Url)
		proxies = append(proxies, reverseProxy)
	}

	return &Backend{
		serverPools: serverPools,
		proxies:     proxies,
		poolsLen:    poolsLen,
		counter:     0,
	}
}
