package repository

import (
	"context"

	"reader-club/internal/domain/entity"

	"github.com/google/uuid"
)

type BookSuggestionRepository interface {
	Save(ctx context.Context, suggestion *entity.BookSuggestion) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.BookSuggestion, error)
	FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.BookSuggestion, error)
	FindEligibleForDraw(ctx context.Context, clubID uuid.UUID) ([]entity.BookSuggestion, error)
	DeleteByClubID(ctx context.Context, clubID uuid.UUID) error
}
