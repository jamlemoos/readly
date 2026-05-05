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

type gormBookSuggestionRepository struct {
	db *gorm.DB
}

func NewBookSuggestionRepository(db *gorm.DB) repository.BookSuggestionRepository {
	return &gormBookSuggestionRepository{db: db}
}

func (r *gormBookSuggestionRepository) Save(ctx context.Context, s *entity.BookSuggestion) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *gormBookSuggestionRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.BookSuggestion, error) {
	var s entity.BookSuggestion
	err := r.db.WithContext(ctx).Preload("User").First(&s, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &s, err
}

func (r *gormBookSuggestionRepository) FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.BookSuggestion, error) {
	var list []entity.BookSuggestion
	err := r.db.WithContext(ctx).Preload("User").Where("club_id = ?", clubID).Find(&list).Error
	return list, err
}

// FindEligibleForDraw returns suggestions not yet selected as a theme.
func (r *gormBookSuggestionRepository) FindEligibleForDraw(ctx context.Context, clubID uuid.UUID) ([]entity.BookSuggestion, error) {
	var list []entity.BookSuggestion
	err := r.db.WithContext(ctx).Preload("User").
		Where("club_id = ? AND id NOT IN (SELECT book_suggestion_id FROM monthly_themes WHERE club_id = ?)", clubID, clubID).
		Find(&list).Error
	return list, err
}

func (r *gormBookSuggestionRepository) DeleteByClubID(ctx context.Context, clubID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("club_id = ?", clubID).Delete(&entity.BookSuggestion{}).Error
}
