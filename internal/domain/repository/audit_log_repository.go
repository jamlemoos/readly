package repository

import (
	"context"

	"reader-club/internal/domain/entity"
)

type AuditLogRepository interface {
	Save(ctx context.Context, log *entity.AuditLog) error
}
