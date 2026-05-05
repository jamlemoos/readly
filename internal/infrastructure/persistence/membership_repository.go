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

type gormMembershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) repository.MembershipRepository {
	return &gormMembershipRepository{db: db}
}

func (r *gormMembershipRepository) Save(ctx context.Context, m *entity.Membership) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *gormMembershipRepository) FindByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (*entity.Membership, error) {
	var m entity.Membership
	err := r.db.WithContext(ctx).Preload("User").
		Where("user_id = ? AND club_id = ?", userID, clubID).
		First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.ErrNotFound
	}
	return &m, err
}

func (r *gormMembershipRepository) FindByClub(ctx context.Context, clubID uuid.UUID) ([]entity.Membership, error) {
	var list []entity.Membership
	err := r.db.WithContext(ctx).Preload("User").Where("club_id = ?", clubID).Find(&list).Error
	return list, err
}

func (r *gormMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]entity.Membership, error) {
	var list []entity.Membership
	err := r.db.WithContext(ctx).Preload("BookClub").Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (r *gormMembershipRepository) CountByClub(ctx context.Context, clubID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Membership{}).Where("club_id = ?", clubID).Count(&count).Error
	return count, err
}

func (r *gormMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Membership{}, "id = ?", id).Error
}

func (r *gormMembershipRepository) DeleteByClubID(ctx context.Context, clubID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("club_id = ?", clubID).Delete(&entity.Membership{}).Error
}

func (r *gormMembershipRepository) ExistsByUserAndClub(ctx context.Context, userID, clubID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Membership{}).
		Where("user_id = ? AND club_id = ?", userID, clubID).
		Count(&count).Error
	return count > 0, err
}
