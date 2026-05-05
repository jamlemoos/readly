package repository

import (
	"context"

	"github.com/google/uuid"
	"reader-club/internal/domain/entity"
)

type BookClubRepository interface {
	Save(ctx context.Context, club *entity.BookClub) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.BookClub, error)
	FindAll(ctx context.Context) ([]entity.BookClub, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ExistsByName(ctx context.Context, name string) (bool, error)
}
