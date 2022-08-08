package ports

import (
	"context"

	"github.com/seggga/approve-mail/internal/domain/models"
)

// SMTPSender sends smtp messages
type SMTPSender interface {
	Send(ctx context.Context, msg *models.MsgSMTP) error
	Notify(ctx context.Context, msg *models.MsgAnalytics) error
}
