package worker

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// remoteConn returns a cached NATS connection to the given domain, or opens a new one if needed.
func (w *Worker) remoteConn(domain string) (*nats.Conn, error) {
	w.remotesMu.Lock()
	defer w.remotesMu.Unlock()

	if re, ok := w.remotes[domain]; ok && re.conn.IsConnected() {
		return re.conn, nil
	}

	wsURL := fmt.Sprintf("ws://%s:%d", domain, w.opts.Port)

	rc, err := nats.Connect(
		wsURL,
		nats.UserInfo("remote", w.opts.RemotePassword),
		nats.Name("simplechat-"+w.opts.Domain+"-remote"),
		nats.MaxReconnects(5),
		nats.ReconnectWait(2*time.Second),
		nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("dial %s (%s): %w", domain, wsURL, err)
	}

	w.remotes[domain] = &remoteEntry{conn: rc}
	log.Printf("[worker] opened remote connection to %s (%s)", domain, wsURL)
	return rc, nil
}

// dropRemote evicts a domain's cached connection.
func (s *Worker) dropRemote(domain string) {
	s.remotesMu.Lock()
	defer s.remotesMu.Unlock()
	if re, ok := s.remotes[domain]; ok {
		_ = re.conn.Drain()
		delete(s.remotes, domain)
	}
}
