package models

// MsgSMTP ...
type MsgSMTP struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Body     []byte `json:"body"`

	EventType string `json:"eventtype"`
	TaskID    uint64 `json:"taskid"`
}
