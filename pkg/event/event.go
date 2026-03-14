package event

import (
	"time"

	"github.com/fibrchat/worker/pkg/address"
)

type Type int

const (
	Connect Type = iota
	Disconnect
)

type Event struct {
	Type      Type            `json:"type"`
	User      address.Address `json:"user"`
	Timestamp time.Time       `json:"ts"`
}
