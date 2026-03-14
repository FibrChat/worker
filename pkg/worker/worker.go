package worker

import (
	"fmt"

	"github.com/fibrchat/server/pkg/subject"
	"github.com/nats-io/nats.go"
)

// Start initializes the worker, connects to NATS, and starts listening for messages.
func Start(o Options) (*Worker, error) {
	if o.Domain == "" {
		return nil, fmt.Errorf("Domain is required")
	}
	if o.ServerURL == "" && o.InProcessServer == nil {
		return nil, fmt.Errorf("ServerURL or InProcessServer is required")
	}
	if o.WorkerPassword == "" {
		return nil, fmt.Errorf("WorkerPassword is required")
	}

	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.Name("worker-" + o.Domain),
		nats.UserInfo("worker", o.WorkerPassword),
	}

	serverURL := o.ServerURL
	if o.InProcessServer != nil {
		opts = append(opts, nats.InProcessServer(o.InProcessServer))
		serverURL = nats.DefaultURL
	}

	nc, err := nats.Connect(serverURL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect NATS: %w", err)
	}

	w := &Worker{nc: nc, opts: o}

	if _, err := nc.QueueSubscribe(subject.PublishSubject, "worker-message", w.handleMessage); err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", subject.PublishSubject, err)
	}

	if _, err := nc.QueueSubscribe(subject.UsersSubject, "worker-users", w.handleUsers); err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", subject.UsersSubject, err)
	}

	if _, err := nc.QueueSubscribe("$SYS.ACCOUNT.*.CONNECT", "worker-presence", w.handleSysConnect); err != nil {
		return nil, fmt.Errorf("subscribe sys connect: %w", err)
	}

	if _, err := nc.QueueSubscribe("$SYS.ACCOUNT.*.DISCONNECT", "worker-presence", w.handleSysDisconnect); err != nil {
		return nil, fmt.Errorf("subscribe sys disconnect: %w", err)
	}

	return w, nil
}

// Stop gracefully stops all components.
func (s *Worker) Stop() {
	s.nc.Drain()
}
