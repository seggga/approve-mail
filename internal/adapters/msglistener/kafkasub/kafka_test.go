// Integration test, depends on running kafka, zookeeper and postgres instances
// started by command: make kafka/compose/up

package kafkasub

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/seggga/approve-mail/internal/domain/models"
	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

const (
	broker  = "127.0.0.1:9093"
	topic   = "test-topic"
	groupID = "test-consumer-group"
)

var (
	c      *Client
	err    error
	chTask chan models.MsgTask

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
	logger, _ := zap.NewDevelopment()
	chTask = make(chan models.MsgTask)
	defer close(chTask)

	c, err = New(broker, topic, groupID, logger, chTask)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	defer c.Reader.Close()

	os.Exit(m.Run())
}

func TestClient(t *testing.T) {
	type publisher struct {
		writer *kafka.Writer
	}

	pub := publisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}

	// convert simple messages to []byte to compose kafka messages
	kfkMessages := make([]kafka.Message, 0, len(taskMessages))
	for _, m := range taskMessages {
		value, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("cannot marshal msg %v, %v", m, err)
		}

		kfkMsg := kafka.Message{
			Value: value,
		}

		kfkMessages = append(kfkMessages, kfkMsg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go c.Start(ctx)
	// send test messages to kafka.
	for i := 0; i < len(kfkMessages); i += 1 {
		err := pub.writer.WriteMessages(context.Background(), kfkMessages[i:i+1]...)
		if err != nil {
			t.Fatal(err)
		}

		select {
		case <-ctx.Done():
			t.Fatalf("exit on timeout")

		case msg := <-chTask:
			if msg != taskMessages[i] {
				t.Fatalf("incorrect data: expected %v, got %v", taskMessages[i], msg)
			}
		}
	}

	pub.writer.Close()
	<-ctx.Done()
}
