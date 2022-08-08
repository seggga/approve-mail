package smtpsender

import (
	"context"
	"time"

	"github.com/seggga/approve-mail/internal/domain/models"
	"github.com/seggga/approve-mail/internal/ports"
	"go.uber.org/zap"
)

var _ ports.SMTPSender = &Sender{}

// Sender ...
type Sender struct {
	chIn     chan models.MsgSMTP
	logger   *zap.Logger
	notifier ports.MsgPubber
	rate     time.Duration
	// smtp-server
}

// New ...
func New(logger *zap.Logger, notifier ports.MsgPubber, ch chan models.MsgSMTP, rateSeconds int) *Sender {
	return &Sender{
		chIn:     ch,
		logger:   logger,
		notifier: notifier,
		rate:     time.Second * time.Duration(rateSeconds),
	}
}

// Start makes seder wait for incoming messages to send them via smtp-server
func (s *Sender) Start(ctx context.Context) error {
	rateLimiter := make(chan struct{}, 1)
	defer close(rateLimiter)

	next := true

	for next {
		select {
		case <-ctx.Done():
			s.logger.Debug("context has been closed")
			next = false
			break

		case msg := <-s.chIn:
			s.logger.Sugar().Debugf("message has been recieved, %v", msg)

			select {
			case <-ctx.Done():
				s.logger.Debug("context has been closed")
				next = false
				break
			case rateLimiter <- struct{}{}:
				err := s.Send(ctx, &msg)
				if err != nil {
					s.logger.Sugar().Errorf("error sending message, %v", err)
					break
				}

				go func() {
					time.Sleep(s.rate)
					<-rateLimiter
				}()

				notifyMsg := models.MsgAnalytics{
					EventType:  msg.EventType,
					TaskID:     msg.TaskID,
					Approver:   msg.Receiver,
					RecievedAt: time.Now(),
				}
				s.logger.Sugar().Debugf("message for analytics has been composed, %v", notifyMsg)

				err = s.Notify(ctx, &notifyMsg)
				if err != nil {
					s.logger.Sugar().Errorf("error sending message with Kafka: %v", err)
					break
				}
				s.logger.Debug("message for analytics has been successfully sent")
			}
		}
	}
	return nil
}

// Send sends messages using smtp server
func (s *Sender) Send(ctx context.Context, msg *models.MsgSMTP) error {
	time.Sleep(time.Second * 5)
	return nil
}

// Notify sends notifications to Analytics microservice
func (s *Sender) Notify(ctx context.Context, msg *models.MsgAnalytics) error {
	return s.notifier.Put(ctx, msg)
}
