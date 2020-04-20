package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type cmdLineOpts struct {
	FlapAfter       time.Duration
	InitWaitSeconds int
}

var (
	opts  cmdLineOpts
	alive bool
	ready bool
)

func init() {
	flag.IntVar(&opts.InitWaitSeconds, "init-wait-seconds", 4, "Time to wait in initilization")
}

func main() {
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("Signal received: ", sig)
		done <- true
	}()

	fmt.Println("Begin application initialization")

	for i := 1; i < opts.InitWaitSeconds; i++ {
		fmt.Printf("Initializing: %d / %d\n", i, opts.InitWaitSeconds)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Initialized")
	alive = true
	fmt.Println("Started but not ready")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've hit the index page\n")
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		if !alive {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("FAILED\n"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\n"))
		}
	})

	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		if !ready {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("FAILED\n"))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK\n"))
		}
	})

	server := &http.Server{Addr: ":8080"}

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
	ready = false
	fmt.Println("We have been signalled to terminate, set the application not ready anymore")
	fmt.Println("Continue to listen for, e.g., 5+ seconds +- till the listener has all connections down")
	time.Sleep(7 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("ERROR: %+v\n", err)
	}
	fmt.Println("Listener has shutdowned")

	fmt.Println("Exiting")
}
