package repository

import (
	"context"

	"reader-club/internal/domain/entity"

	"github.com/google/uuid"
)

type MeetingRepository interface {
	Save(ctx context.Context, meeting *entity.Meeting) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error)
	FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.Meeting, error)
	DeleteByClubID(ctx context.Context, clubID uuid.UUID) error
}
