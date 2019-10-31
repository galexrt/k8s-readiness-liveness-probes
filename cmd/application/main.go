package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	alive bool
	ready bool
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("Begin application initialization")

	waitTime := 4
	for i := 1; i < waitTime; i++ {
		fmt.Printf("Initializing:%d / %d\n", i, waitTime)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Initialized")
	alive = true
	fmt.Println("Started but not ready")

	server := &http.Server{Addr: ":8080", Handler: http.HandlerFunc(serve)}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("ERROR: %+v\n", err)
		}
	}()

	fmt.Println("We need 3 more seconds before we are ready (e.g., some more db connections or so")
	time.Sleep(3 * time.Second)
	ready = true

	fmt.Println("Ready, and waiting for signal")
	<-done
	fmt.Println("We should terminate, not ready anymore")
	fmt.Println("Continue to listen for, e.g., 5+ seconds +- till the listener has all connections down")
	ready = false
	time.Sleep(7 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
	fmt.Println("Last connections done, listener has shutdowned")

	fmt.Println("Exiting")
}

func serve(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
	case "/liveness":
		if !alive {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("FAILED\n"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\n"))
		}
	case "/readiness":
		if !ready {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("FAILED\n"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\n"))
		}
	}
}
