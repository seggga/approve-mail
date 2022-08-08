package kafkapub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/seggga/approve-mail/internal/domain/models"
	"github.com/seggga/approve-mail/internal/ports"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var _ ports.MsgPubber = &Client{}

// Client ...
type Client struct {
	writer *kafka.Writer
	logger *zap.Logger
}

// New ...
func New(broker, topic string, logger *zap.Logger) (*Client, error) {
	if broker == "" || topic == "" {
		return nil, fmt.Errorf("missed some connection parameters: brokers %s, topic %s", broker, topic)
	}

	c := Client{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
		logger: logger,
	}

	return &c, nil
}

// Put writes messages in Kafka for Analytics microservice
func (c *Client) Put(ctx context.Context, msg *models.MsgAnalytics) error {

	value, err := json.Marshal(*msg)
	if err != nil {
		c.logger.Sugar().Errorf("cannot marshal msg %v, %v", msg, err)
		return fmt.Errorf("cannot marshal message: %w", err)
	}

	var messages []kafka.Message
	kfkMsg := kafka.Message{
		Value: value,
	}
	messages = append(messages, kfkMsg)
	err = c.writer.WriteMessages(ctx, messages...)
	if err != nil {
		c.logger.Sugar().Errorf("cannot send message to kafka %v, %v", msg, err)
		return fmt.Errorf("cannot send message to kafka: %w", err)
	}

	return nil
}
