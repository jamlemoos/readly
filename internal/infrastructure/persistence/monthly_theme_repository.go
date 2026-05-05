package persistence

import (
	"context"
	"errors"

	"reader-club/internal/application/apperr"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type gormMonthlyThemeRepository struct {
	db *gorm.DB
}

func NewMonthlyThemeRepository(db *gorm.DB) repository.MonthlyThemeRepository {
	return &gormMonthlyThemeRepository{db: db}
}

func (r *gormMonthlyThemeRepository) Save(ctx context.Context, t *entity.MonthlyTheme) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *gormMonthlyThemeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.MonthlyTheme, error) {
	var t entity.MonthlyTheme
	err := r.db.WithContext(ctx).Preload("BookSuggestion.User").First(&t, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &t, err
}

func (r *gormMonthlyThemeRepository) FindByClubAndMonth(ctx context.Context, clubID uuid.UUID, year, month int) (*entity.MonthlyTheme, error) {
	var t entity.MonthlyTheme
	err := r.db.WithContext(ctx).Preload("BookSuggestion.User").
		Where("club_id = ? AND year = ? AND month = ?", clubID, year, month).
		First(&t).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &t, err
}

func (r *gormMonthlyThemeRepository) FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.MonthlyTheme, error) {
	var list []entity.MonthlyTheme
	err := r.db.WithContext(ctx).Preload("BookSuggestion.User").
		Where("club_id = ?", clubID).Order("year DESC, month DESC").Find(&list).Error
	return list, err
}

func (r *gormMonthlyThemeRepository) ExistsByClubAndMonth(ctx context.Context, clubID uuid.UUID, year, month int) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.MonthlyTheme{}).
		Where("club_id = ? AND year = ? AND month = ?", clubID, year, month).
		Count(&count).Error
	return count > 0, err
}

func (r *gormMonthlyThemeRepository) DeleteByClubID(ctx context.Context, clubID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("club_id = ?", clubID).Delete(&entity.MonthlyTheme{}).Error
}
