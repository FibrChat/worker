package worker

import (
	"encoding/json"
	"log"

	"github.com/fibrchat/worker/pkg/address"
	"github.com/fibrchat/worker/pkg/message"

	"github.com/nats-io/nats.go"
)

// handleRemote receives messages from remote clients and routes them to the appropriate local recipient.
func (s *Worker) handleRemote(msg *nats.Msg) {
	var cm message.Message
	if err := json.Unmarshal(msg.Data, &cm); err != nil {
		log.Printf("[worker] malformed remote message: %v", err)
		s.respond(msg, message.CodeInternal, "malformed message")
		return
	}

	if err := s.sendLocal(cm); err != nil {
		log.Printf("[worker] sendLocal error: %v", err)
		s.respond(msg, message.CodeInternal, err.Error())
		return
	}

	s.respond(msg, message.CodeSuccess, "")
}

// handleLocal receives messages from local clients and routes them to the appropriate destination (local or remote).
func (s *Worker) handleLocal(msg *nats.Msg) {
	var cm message.Message
	if err := json.Unmarshal(msg.Data, &cm); err != nil {
		log.Printf("[worker] malformed chat.send: %v", err)
		s.respond(msg, message.CodeInternal, "malformed message")
		return
	}

	if err := s.routeMessage(cm); err != nil {
		log.Printf("[worker] routeMessage error: %v", err)
		s.respond(msg, message.CodeInternal, err.Error())
		return
	}

	s.respond(msg, message.CodeSuccess, "")
}

// routeMessage determines whether the message is destined for a local or remote recipient and routes it accordingly.
func (s *Worker) routeMessage(cm message.Message) error {
	_, domain, err := address.Split(cm.To)
	if err != nil {
		return err
	}

	if domain == s.opts.Domain {
		return s.sendLocal(cm)
	}

	return s.sendRemote(cm, domain)
}
