package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         uuid.UUID       `gorm:"type:uuid;primaryKey"`
	ActorID    uuid.UUID       `gorm:"type:uuid;not null"`
	Action     string          `gorm:"not null"`
	EntityType string          `gorm:"not null"`
	EntityID   uuid.UUID       `gorm:"type:uuid"`
	Metadata   json.RawMessage `gorm:"type:jsonb"`
	CreatedAt  time.Time

	Actor User `gorm:"foreignKey:ActorID"`
}

func NewAuditLog(actorID uuid.UUID, action, entityType string, entityID uuid.UUID, metadata any) (*AuditLog, error) {
	raw, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	return &AuditLog{
		ID:         uuid.New(),
		ActorID:    actorID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Metadata:   raw,
	}, nil
}
