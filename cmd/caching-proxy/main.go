package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	proxy "github.com/duanechan/caching-proxy/internal"
)

func main() {
	origin := flag.String("origin", "", "the url address to forward requests to and cache from")
	port := flag.Int("port", 3000, "the port number of the proxy server")
	flag.Parse()

	proxyHandler, err := proxy.NewProxy(*origin, *port)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler)
	address := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(address, mux))
}
