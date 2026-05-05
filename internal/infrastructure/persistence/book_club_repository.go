package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"reader-club/internal/application/apperr"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type gormBookClubRepository struct {
	db *gorm.DB
}

func NewBookClubRepository(db *gorm.DB) repository.BookClubRepository {
	return &gormBookClubRepository{db: db}
}

func (r *gormBookClubRepository) Save(ctx context.Context, club *entity.BookClub) error {
	return r.db.WithContext(ctx).Save(club).Error
}

func (r *gormBookClubRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.BookClub, error) {
	var club entity.BookClub
	err := r.db.WithContext(ctx).Preload("Owner").First(&club, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &club, err
}

func (r *gormBookClubRepository) FindAll(ctx context.Context) ([]entity.BookClub, error) {
	var clubs []entity.BookClub
	err := r.db.WithContext(ctx).Find(&clubs).Error
	return clubs, err
}

func (r *gormBookClubRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&entity.BookClub{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *gormBookClubRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.BookClub{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}
