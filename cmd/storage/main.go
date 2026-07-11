package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/SteGG200/storage/internal/server"
)

func main() {
	var portFlag string
	flag.StringVar(&portFlag, "port", "", "port to listen on (e.g. 8080)")
	flag.StringVar(&portFlag, "p", "", "port to listen on (e.g. 8080) (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] /path/to/storage\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	storagePath := args[0]
	absPath, err := filepath.Abs(storagePath)
	if err != nil {
		log.Fatalf("Error determining absolute path: %v", err)
	}

	// #nosec
	fi, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			// #nosec G706
			log.Fatalf("Error: storage path %q does not exist", absPath)
		}
		log.Fatalf("Error checking storage path: %v", err)
	}

	if !fi.IsDir() {
		// #nosec G706
		log.Fatalf("Error: storage path %q is not a directory", absPath)
	}

	port := portFlag
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	cfg := server.Config{
		Addr:        addr,
		StorageRoot: absPath,
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	srvErr := make(chan error, 1)
	go func() {
		// #nosec G706
		log.Printf("Starting storage API server on %s, serving root: %s", addr, absPath)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		log.Fatalf("Server error: %v", err)
	case sig := <-shutdown:
		log.Printf("Received signal %v, shutting down gracefully...", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			_ = srv.Close()
		}
		log.Println("Server stopped.")
	}
}
