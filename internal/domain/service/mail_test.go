package service

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/seggga/approve-mail/internal/domain/models"
)

const ()

var (
	srv   *Service
	chIn  chan models.MsgTask
	chOut chan models.MsgSMTP

	taskMessages []models.MsgTask = []models.MsgTask{
		// approved task
		{
			TaskID:      1,
			TaskName:    "approve task 1",
			TaskDescr:   "approve task description 1",
			EventType:   models.Approval,
			Approver:    "one@mail.com",
			AcceptLink:  "http://accept",
			DeclineLink: "http://decline",
		},

		{
			TaskID:    2,
			TaskName:  "decline task 2",
			TaskDescr: "decline task description 2",
			EventType: models.Decline,
			Approver:  "two@mail.com",
		},
		{
			TaskID:      3,
			TaskName:    "incoorrect task 3",
			TaskDescr:   "approve task description 3",
			EventType:   "incorrect_event_type",
			Approver:    "three@mail.com",
			AcceptLink:  "http://accept",
			DeclineLink: "http://decline",
		},
	}
)

func TestMain(m *testing.M) {
	chIn = make(chan models.MsgTask)
	chOut = make(chan models.MsgSMTP)
	srv, _ = New(chIn, chOut, 2)
	ctx := context.TODO()
	go srv.Start(ctx)

	os.Exit(m.Run())
}

// expect compose approve message
func TestComposeApprove(t *testing.T) {
	chIn <- taskMessages[0]
	out := <-chOut

	expect := models.MsgSMTP{
		Sender:   sender,
		Receiver: taskMessages[0].Approver,
	}
	var body bytes.Buffer
	approveTmpl.Execute(&body, taskMessages[0])
	expect.Body = body.Bytes()

	if out.Sender != expect.Sender {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", expect.Sender, out.Sender)
	}
	if out.Receiver != expect.Receiver {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", expect.Receiver, out.Receiver)
	}
	if string(out.Body) != string(expect.Body) {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", string(expect.Body), string(out.Body))
	}
}

// expect compose decline message
func TestComposeDecline(t *testing.T) {
	chIn <- taskMessages[1]
	out := <-chOut

	expect := models.MsgSMTP{
		Sender:   sender,
		Receiver: taskMessages[1].Approver,
	}
	var body bytes.Buffer
	declTmpl.Execute(&body, taskMessages[1])
	expect.Body = body.Bytes()

	if out.Sender != expect.Sender {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", expect.Sender, out.Sender)
	}
	if out.Receiver != expect.Receiver {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", expect.Receiver, out.Receiver)
	}
	if string(out.Body) != string(expect.Body) {
		t.Errorf("Approve set: incorrect sender: expected %s, got %s", string(expect.Body), string(out.Body))
	}
}

// expect no messages in output channel
func TestComposeIncorrect(t *testing.T) {
	chIn <- taskMessages[2]
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		break
	case out := <-chOut:
		t.Errorf("unexpected output: %v", out)
	}
}
