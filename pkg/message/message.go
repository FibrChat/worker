package message

import "time"

type Message struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"ts"`
}

type Response struct {
	Code  int    `json:"code"`
	Error string `json:"error,omitempty"`
}

const (
	CodeSuccess = iota
	CodeInternal
)
