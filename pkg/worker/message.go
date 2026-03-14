package worker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fibrchat/server/pkg/subject"
	"github.com/fibrchat/worker/pkg/message"
	"github.com/fibrchat/worker/pkg/request"
	"github.com/nats-io/nats.go"
)

// sendMessage publishes the message to the recipient's DM NATS subject.
func (w *Worker) sendMessage(cm message.Message) error {
	data, _ := json.Marshal(cm)
	if err := w.nc.Publish(subject.Inbox(cm.Dst.ID), data); err != nil {
		log.Printf("[worker] deliver to %q: %v", cm.Dst, err)
		return fmt.Errorf("failed to deliver to %q: %w", cm.Dst, err)
	}

	return nil
}

// sendReply sends a Reponse back to the client
func (w *Worker) sendReply(msg *nats.Msg, code request.ResponseCode, errMsg string) {
	if msg.Reply == "" {
		return
	}

	reply := request.Response{Code: code}
	if errMsg != "" {
		reply.Error = errMsg
	}

	data, _ := json.Marshal(reply)
	if err := msg.Respond(data); err != nil {
		log.Printf("[worker] respond: %v", err)
	}
}
