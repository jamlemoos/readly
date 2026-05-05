package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookClub struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"not null"`
	Description string
	OwnerID     uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Owner       User             `gorm:"foreignKey:OwnerID"`
	Memberships []Membership     `gorm:"foreignKey:ClubID"`
	Suggestions []BookSuggestion `gorm:"foreignKey:ClubID"`
	Themes      []MonthlyTheme   `gorm:"foreignKey:ClubID"`
	Meetings    []Meeting        `gorm:"foreignKey:ClubID"`
}

const MaxMembers = 20

func NewBookClub(name, description string, ownerID uuid.UUID) *BookClub {
	return &BookClub{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
}

func (c *BookClub) CanAddMember(currentCount int) bool {
	return currentCount < MaxMembers
}
