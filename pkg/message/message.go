package message

import (
	"time"

	"github.com/fibrchat/worker/pkg/address"
)

type MessageType int

const (
	MessageDM MessageType = iota
	MessageROOM
)

type Message struct {
	Type MessageType `json:"type"`

	Src address.Address `json:"src"`
	Dst address.Address `json:"dst"`

	Content   string    `json:"content"`
	Timestamp time.Time `json:"ts"`
}

