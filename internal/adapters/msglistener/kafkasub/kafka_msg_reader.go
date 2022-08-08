package kafkasub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seggga/approve-mail/internal/domain/models"
	"github.com/seggga/approve-mail/internal/ports"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var _ ports.MsgSubber = &Client{}

// Client ...
type Client struct {
	Reader *kafka.Reader
	logger *zap.Logger
	chTask chan models.MsgTask
}

// New ...
func New(broker, topic, groupID string, logger *zap.Logger, ch chan models.MsgTask) (*Client, error) {
	if broker == "" || topic == "" || groupID == "" {
		return nil, fmt.Errorf("missed some connection parameters: brokers %s, topic %s, groupID %s", broker, topic, groupID)
	}

	c := Client{
		logger: logger,
		chTask: ch,
	}

	c.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e1,
		MaxBytes: 10e6,
	})

	return &c, nil
}

// Get calls (ports.Analyter).WriteEvent method that processes incoming message
func (c *Client) Get(ctx context.Context, msg *models.MsgTask) error {

	select {
	case c.chTask <- *msg:
		c.logger.Debug("message has been successfully sent through the channel")

	case <-ctx.Done():
		c.logger.Debug("context has been cancelled")
	}

	return nil

}

// Start sequentaly reads messages from kafka, commits messages after processing
func (c *Client) Start(ctx context.Context) error {
	msg := &models.MsgTask{}
	next := true

	for next {
		select {

		case <-ctx.Done():
			c.logger.Debug("context has been closed")
			next = false
			break

		default:
			kafkaMsg, err := c.Reader.FetchMessage(ctx)
			if err != nil {
				c.logger.Sugar().Debugf("error fetching message: %v", err)
				continue
			}

			c.logger.Sugar().Debugf("got []byte message %v", kafkaMsg)

			err = json.Unmarshal(kafkaMsg.Value, msg)
			if err != nil {
				c.logger.Sugar().Debugf("failed unmarshal kafka message %v: %v", kafkaMsg, err)
				continue
			}

			c.logger.Sugar().Debugf("parsed message %v", msg)

			if err := c.Get(ctx, msg); err != nil {
				c.logger.Sugar().Debugf("failed message processing %v: %v", msg, err)
				continue
			}

			if err := c.Reader.CommitMessages(ctx, kafkaMsg); err != nil {
				c.logger.Sugar().Debugf("error committing message: %v", err)
				continue
			}
		}
	}
	return nil
}

// Stop ...
func (c *Client) Stop() error {
	return c.Reader.Close()
}
