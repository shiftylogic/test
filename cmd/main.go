package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	port = ":9337"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/readiness", readyHandler)
	mux.HandleFunc("/", allTheThings)

	// Setup the server structures
	server := &http.Server{
		Handler:      mux,
		Addr:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	idleClosed := make(chan struct{})

	// Clean shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown: %v\n", err)
		}

		close(idleClosed)
	}()

	// Start the service
	log.Printf("Starting server (port: %s)\n", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	<-idleClosed
	log.Println("Goodbye.")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func allTheThings(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, from LivingLogical\n\nPath: %q", html.EscapeString(r.URL.Path))
}
