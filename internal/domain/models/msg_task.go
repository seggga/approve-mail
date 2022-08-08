package models

const (
	// Approval is one of expected MsgTask.EventType values
	// means we have to use approval message template
	Approval = "approval"

	// Decline is one of expected MsgTask.EventType values
	// means service will use decline message template
	Decline = "decine"
)

// MsgTask represents a message from Tasks service
// It is expected to recieve two types of messages:
//   - ask for approve message
//   - message on declined task
// Type of message is defined by kafka.Key
type MsgTask struct {
	TaskID      uint64 // task ID
	TaskName    string // task name
	TaskDescr   string // task description
	EventType   string // event type (send approve message / decline)
	Approver    string // recipient's smtp-address
	AcceptLink  string // a link to accept the task for specified approver
	DeclineLink string // a link to decline given task
}
