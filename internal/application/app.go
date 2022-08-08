package application

import (
	"context"

	"github.com/seggga/approve-mail/internal/adapters/msglistener/kafkasub"
	"github.com/seggga/approve-mail/internal/adapters/msgwriter/kafkapub"
	"github.com/seggga/approve-mail/internal/adapters/smtpsender"
	"github.com/seggga/approve-mail/internal/domain/models"
	"github.com/seggga/approve-mail/internal/domain/service"
	"golang.org/x/sync/errgroup"

	"go.uber.org/zap"
)

var (
	kafkaSub   *kafkasub.Client
	smtpSender *smtpsender.Sender
	composer   *service.Service
	logger     *zap.Logger

	chanTask chan models.MsgTask
	chanSMTP chan models.MsgSMTP
)

// Start ...
func Start(ctx context.Context) {
	chanTask = make(chan models.MsgTask)
	chanSMTP = make(chan models.MsgSMTP)

	cfg := getConfig()
	logger := initLogger(cfg.Logger.Level)

	var err error
	kafkaSub, err = kafkasub.New(cfg.Kafka.Sub.Server, cfg.Kafka.Sub.Topic, cfg.Kafka.Sub.GroupID, logger, chanTask)
	if err != nil {
		logger.Sugar().Fatalf("cannot create Kafka reader: %w", err)
	}

	kafkaPub, err := kafkapub.New(cfg.Kafka.Pub.Server, cfg.Kafka.Pub.Topic, logger)
	if err != nil {
		logger.Sugar().Fatalf("cannot create Kafka publisher: %w", err)
	}
	smtpSender = smtpsender.New(logger, kafkaPub, chanSMTP, cfg.Mail.Rate)

	composer, err = service.New(chanTask, chanSMTP, cfg.Compose.Workers)

	var g errgroup.Group
	g.Go(func() error {
		return kafkaSub.Start(ctx)
	})
	g.Go(func() error {
		return smtpSender.Start(ctx)
	})
	g.Go(func() error {
		return composer.Start(ctx)
	})

	logger.Info("app is started")
	err = g.Wait()
	if err != nil {
		logger.Sugar().Fatalf("applicateion failed: %v", err)
	}

}

// Stop ...
func Stop() {
	defer logger.Sync()
	defer close(chanSMTP)
	defer close(chanTask)

	// stop reading messages from task microservice via kafka
	err := kafkaSub.Stop()
	if err != nil {
		logger.Sugar().Errorf("error stopping kafka reader", err)
	}

	logger.Sugar().Info("application has been stopped")

}
