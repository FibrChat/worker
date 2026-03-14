package worker

import (
	"encoding/json"
	"log"

	"github.com/fibrchat/worker/pkg/message"

	"github.com/nats-io/nats.go"
)

// handleMessage handles messages
func (w *Worker) handleMessage(msg *nats.Msg) {
	var cm message.Message
	if err := json.Unmarshal(msg.Data, &cm); err != nil {
		log.Printf("[worker] malformed chat.send: %v", err)
		w.sendReply(msg, message.CodeInternal, "malformed message")
		return
	}

	if err := w.sendMessage(cm); err != nil {
		log.Printf("[worker] sendMessage error: %v", err)
		w.sendReply(msg, message.CodeInternal, err.Error())
		return
	}

	w.sendReply(msg, message.CodeSuccess, "")
}
