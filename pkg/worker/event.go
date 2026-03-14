package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/fibrchat/server/pkg/subject"
	"github.com/fibrchat/worker/pkg/address"
	"github.com/fibrchat/worker/pkg/event"
	"github.com/nats-io/nats.go"
)

type sysConnectEvent struct {
	Client struct {
		User string `json:"user"`
	} `json:"client"`
}

func (w *Worker) handleSysConnect(msg *nats.Msg) {
	w.publishPresence(msg, event.Connect)
}

func (w *Worker) handleSysDisconnect(msg *nats.Msg) {
	w.publishPresence(msg, event.Disconnect)
}

func (w *Worker) publishPresence(msg *nats.Msg, typ event.Type) {
	var sys sysConnectEvent
	if err := json.Unmarshal(msg.Data, &sys); err != nil {
		return
	}

	if sys.Client.User == "" || sys.Client.User == "worker" {
		return
	}

	evt := event.Event{
		Type:      typ,
		User:      address.Address{ID: sys.Client.User, Domain: w.opts.Domain},
		Timestamp: time.Now().UTC(),
	}

	data, _ := json.Marshal(evt)
	if err := w.nc.Publish(subject.PresenceSubject, data); err != nil {
		log.Printf("[worker] publish presence: %v", err)
	}
}
