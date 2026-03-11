package worker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fibrchat/server/pkg/subject"
	"github.com/fibrchat/worker/pkg/address"
	"github.com/fibrchat/worker/pkg/message"
	"github.com/nats-io/nats.go"
)

// sendLocal publishes the message to the recipient's DM NATS subject.
func (s *Worker) sendLocal(cm message.Message) error {
	name, _, err := address.Split(cm.To)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(cm)
	if err := s.nc.Publish(subject.DM(name), data); err != nil {
		log.Printf("[worker] deliver to %q: %v", name, err)
		return fmt.Errorf("failed to deliver to %q: %w", name, err)
	}

	return nil
}

// sendRemote publishes the message to the remote domain's NATS subject.
func (s *Worker) sendRemote(cm message.Message, domain string) error {
	rc, err := s.remoteConn(domain)
	if err != nil {
		return fmt.Errorf("cannot reach %s: %w", domain, err)
	}

	data, err := json.Marshal(cm)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if err := rc.Publish(subject.Remote, data); err != nil {
		s.dropRemote(domain)
		return fmt.Errorf("failed to send to %s: %w", domain, err)
	}

	log.Printf("[worker] sent remote message to %s", cm.To)

	return nil
}

// respond sends a Reponse back to the client
func (s *Worker) respond(msg *nats.Msg, code int, errMsg string) {
	if msg.Reply == "" {
		return
	}

	reply := message.Response{Code: code}
	if errMsg != "" {
		reply.Error = errMsg
	}

	data, _ := json.Marshal(reply)
	if err := msg.Respond(data); err != nil {
		log.Printf("[worker] respond: %v", err)
	}
}
