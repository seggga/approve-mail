package models

import "time"

// MsgAnalytics represents outgoing message to Analytics service
type MsgAnalytics struct {
	EventType  string    `json:"eventtype"`
	TaskID     uint64    `json:"taskid"`
	Approver   string    `json:"approver"`
	RecievedAt time.Time `json:"recievedat"`
}
