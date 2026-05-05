package repository

import (
	"context"

	"reader-club/internal/domain/entity"

	"github.com/google/uuid"
)

type MonthlyThemeRepository interface {
	Save(ctx context.Context, theme *entity.MonthlyTheme) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.MonthlyTheme, error)
	FindByClubAndMonth(ctx context.Context, clubID uuid.UUID, year, month int) (*entity.MonthlyTheme, error)
	FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.MonthlyTheme, error)
	ExistsByClubAndMonth(ctx context.Context, clubID uuid.UUID, year, month int) (bool, error)
	DeleteByClubID(ctx context.Context, clubID uuid.UUID) error
}
