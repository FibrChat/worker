package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// Start initializes the worker, connects to NATS, and starts listening for messages.
func Start(o Options) (*Worker, error) {
	if o.Domain == "" {
		return nil, fmt.Errorf("Domain is required")
	}
	if o.ServerURL == "" {
		return nil, fmt.Errorf("ServerURL is required")
	}
	if o.WorkerPassword == "" {
		return nil, fmt.Errorf("WorkerPassword is required")
	}
	if o.RemotePassword == "" {
		return nil, fmt.Errorf("RemotePassword is required")
	}

	nc, err := nats.Connect(
		o.ServerURL,
		nats.UserInfo("worker", o.WorkerPassword),
		nats.Name("simplechat-worker-"+o.Domain),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect NATS at %s: %w", o.ServerURL, err)
	}

	w := &Worker{
		nc:      nc,
		opts:    o,
		remotes: make(map[string]*remoteEntry),
	}

	if err := w.listen(); err != nil {
		log.Fatalf("Failed to subscribe to NATS subjects: %v", err)
	}

	return w, nil
}

// Shutdown gracefully stops all components.
func (s *Worker) Shutdown() {
	s.nc.Drain()
	s.remotesMu.Lock()
	for _, re := range s.remotes {
		_ = re.conn.Drain()
	}
	s.remotesMu.Unlock()
}
