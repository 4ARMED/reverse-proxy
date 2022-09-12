package main

import (
	"flag"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Largely taken from https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	if targetHost == "" {
		return nil, fmt.Errorf("empty target URL")
	}
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) gin.HandlerFunc {
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	targetHost := flag.String("target", "", "target URL, for example http://loki:3000")
	port := flag.Int("port", 3000, "port to listen on")
	flag.Parse()

	proxy, err := NewProxy(*targetHost)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf(":%d", *port)

	r := gin.Default()
	r.NoRoute(ProxyRequestHandler(proxy))

	log.Printf("starting proxy on %s\n", addr)
	r.Run(addr)
}
