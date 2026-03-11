package worker

import (
	"sync"

	nats "github.com/nats-io/nats.go"
)

type remoteEntry struct {
	conn *nats.Conn
}

type Worker struct {
	opts      Options
	nc        *nats.Conn
	remotesMu sync.Mutex
	remotes   map[string]*remoteEntry
}

type Options struct {
	Port           int
	Domain         string
	ServerURL      string
	WorkerPassword string
	RemotePassword string
}
