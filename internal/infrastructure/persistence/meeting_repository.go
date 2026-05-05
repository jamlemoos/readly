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

type gormMeetingRepository struct {
	db *gorm.DB
}

func NewMeetingRepository(db *gorm.DB) repository.MeetingRepository {
	return &gormMeetingRepository{db: db}
}

func (r *gormMeetingRepository) Save(ctx context.Context, m *entity.Meeting) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *gormMeetingRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Meeting, error) {
	var m entity.Meeting
	err := r.db.WithContext(ctx).
		Preload("Theme.BookSuggestion.User").
		First(&m, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &m, err
}

func (r *gormMeetingRepository) FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.Meeting, error) {
	var list []entity.Meeting
	err := r.db.WithContext(ctx).
		Preload("Theme.BookSuggestion.User").
		Where("club_id = ?", clubID).Order("scheduled_at ASC").Find(&list).Error
	return list, err
}

func (r *gormMeetingRepository) DeleteByClubID(ctx context.Context, clubID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("club_id = ?", clubID).Delete(&entity.Meeting{}).Error
}
