package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	proxy "github.com/duanechan/caching-proxy/internal"
)

func main() {
	origin := flag.String("origin", "", "the origin server to forward requests to and cache responses from")
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
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	go func() {
		proxy.ServerLog("Listening on port", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			proxy.ErrorLog("Error:", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Print("\r")
	proxy.WarnLog("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		proxy.ErrorLog("Shutdown error:", err.Error())
	} else {
		proxy.ServerLog("Server stopped gracefully.")
	}
}
