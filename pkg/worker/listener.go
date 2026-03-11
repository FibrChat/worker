package worker

import (
	"fmt"
	"github.com/fibrchat/server/pkg/subject"
)

// listen subscribes to the NATS subjects for local and remote messages.
func (s *Worker) listen() error {
	if _, err := s.nc.QueueSubscribe(subject.Send, "worker-send", s.handleLocal); err != nil {
		return fmt.Errorf("subscribe %s: %w", subject.Send, err)
	}

	if _, err := s.nc.QueueSubscribe(subject.Remote, "worker-remote", s.handleRemote); err != nil {
		return fmt.Errorf("subscribe %s: %w", subject.Remote, err)
	}

	return nil
}
