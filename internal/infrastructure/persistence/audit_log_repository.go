package persistence

import (
	"context"

	"gorm.io/gorm"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type gormAuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &gormAuditLogRepository{db: db}
}

func (r *gormAuditLogRepository) Save(ctx context.Context, log *entity.AuditLog) error {
	return r.db.WithContext(ctx).Save(log).Error
}
