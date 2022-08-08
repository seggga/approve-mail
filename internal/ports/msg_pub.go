package ports

import (
	"context"

	"github.com/seggga/approve-mail/internal/domain/models"
)

// MsgPubber writes messages for analytics microservice
type MsgPubber interface {
	Put(ctx context.Context, msg *models.MsgAnalytics) error
}
