package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/fibrchat/worker/pkg/request"
	"github.com/nats-io/nats.go"
)

type connzResponse struct {
	Data struct {
		Connections []struct {
			AuthorizedUser string `json:"authorized_user"`
		} `json:"connections"`
	} `json:"data"`
}

func (w *Worker) handleUsers(msg *nats.Msg) {
	resp, err := w.nc.Request("$SYS.REQ.SERVER.PING.CONNZ", []byte(`{"auth":true}`), 2*time.Second)
	if err != nil {
		log.Printf("[worker] connz request: %v", err)
		return
	}

	var connz connzResponse
	if err := json.Unmarshal(resp.Data, &connz); err != nil {
		log.Printf("[worker] connz unmarshal: %v", err)
		return
	}

	var users []string
	for _, c := range connz.Data.Connections {
		if c.AuthorizedUser == "" || c.AuthorizedUser == "worker" {
			continue
		}
		users = append(users, c.AuthorizedUser+"@"+w.opts.Domain)
	}

	reply := request.UsersResponse{Users: users}
	data, _ := json.Marshal(reply)
	if err := msg.Respond(data); err != nil {
		log.Printf("[worker] respond users: %v", err)
	}
}
