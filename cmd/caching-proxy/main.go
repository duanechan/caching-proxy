package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	proxy "github.com/duanechan/caching-proxy/internal"
)

func main() {
	origin := flag.String("origin", "", "the url address to forward requests to and cache from")
	port := flag.Int("port", 3000, "the port number of the proxy server")
	clearCache := flag.Bool("clear-cache", false, "option to clear all cache")
	flag.Parse()

	client := proxy.NewClient(5 * time.Second)
	if *clearCache {
		if err := client.FlushCache(); err != nil {
			proxy.ErrorLog("Error:", err.Error())
		}
		os.Exit(0)
	}

	proxyHandler, err := proxy.NewProxy(client, *origin, *port)
	if err != nil {
		proxy.ErrorLog("Error:", err.Error())
	}

	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler)
	address := fmt.Sprintf(":%d", *port)

	go func() {
		proxy.ServerLog("Listening on port", address)
		if err := http.ListenAndServe(address, mux); err != nil {
			proxy.ErrorLog("Error:", err.Error())
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	fmt.Print("\r")
	proxy.ServerLog("Shutting down...")
}
