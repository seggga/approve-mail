package ports

import (
	"context"

	"github.com/seggga/approve-mail/internal/domain/models"
)

// MsgSubber listens for messages from task microservice
type MsgSubber interface {
	Get(ctx context.Context, msg *models.MsgTask) error
}
