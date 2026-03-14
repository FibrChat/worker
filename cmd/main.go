package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fibrchat/worker/pkg/worker"
)

func main() {
	domain := mustEnv("DOMAIN")
	serverURL := mustEnv("SERVER_URL")
	workerPassword := mustEnv("WORKER_PASSWORD")

	wrk, err := worker.Start(worker.Options{
		Domain:         domain,
		ServerURL:      serverURL,
		WorkerPassword: workerPassword,
	})
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	fmt.Printf("Worker running — domain=%s  nats=%s\n", domain, serverURL)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")

	wrk.Stop()
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return v
}
