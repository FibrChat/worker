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

type ResponseCode int

const (
	CodeSuccess ResponseCode = iota
	CodeInternal
)

type Message struct {
	Type MessageType `json:"type"`

	Src address.Address `json:"src"`
	Dst address.Address `json:"dst"`

	Content   string    `json:"content"`
	Timestamp time.Time `json:"ts"`
}

type Response struct {
	Code  ResponseCode `json:"code"`
	Error string       `json:"error,omitempty"`
}
