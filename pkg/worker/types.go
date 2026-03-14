package worker

import (
	nats "github.com/nats-io/nats.go"
)

type Worker struct {
	opts Options
	nc   *nats.Conn
}

type Options struct {
	Domain          string
	ServerURL       string
	WorkerPassword  string
	InProcessServer nats.InProcessConnProvider
}
