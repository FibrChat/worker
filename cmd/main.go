package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fibrchat/worker/pkg/worker"
)

func main() {
	domain := mustEnv("DOMAIN")
	serverURL := mustEnv("SERVER_URL")
	workerPassword := mustEnv("WORKER_PASSWORD")
	remotePassword := envOr("REMOTE_PASSWORD", "simplechat-remote")
	port := envInt("PORT", 4222)

	wrk, err := worker.New(worker.Options{
		Domain:         domain,
		ServerURL:      serverURL,
		WorkerPassword: workerPassword,
		RemotePassword: remotePassword,
		Port:           port,
	})
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	fmt.Printf("SimpleChat worker running — domain=%s  nats=%s\n", domain, serverURL)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")

	wrk.Shutdown()
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid value for %s: %v", key, err)
	}
	return n
}
