package repository

import (
	"context"

	"reader-club/internal/domain/entity"

	"github.com/google/uuid"
)

type MembershipRepository interface {
	Save(ctx context.Context, membership *entity.Membership) error
	FindByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*entity.Membership, error)
	FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.Membership, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]entity.Membership, error)
	CountByClub(ctx context.Context, clubID uuid.UUID) (int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByClubID(ctx context.Context, clubID uuid.UUID) error
	ExistsByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (bool, error)
}
