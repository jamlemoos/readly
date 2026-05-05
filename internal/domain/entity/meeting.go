package entity

import (
	"time"

	"github.com/google/uuid"
)

type Meeting struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClubID      uuid.UUID `gorm:"type:uuid;not null"`
	ThemeID     uuid.UUID `gorm:"type:uuid"`
	ScheduledAt time.Time `gorm:"not null"`
	Location    string
	Notes       string
	CreatedAt   time.Time

	Club  BookClub     `gorm:"foreignKey:ClubID"`
	Theme MonthlyTheme `gorm:"foreignKey:ThemeID"`
}

func NewMeeting(clubID, themeID uuid.UUID, scheduledAt time.Time, location, notes string) *Meeting {
	return &Meeting{
		ID:          uuid.New(),
		ClubID:      clubID,
		ThemeID:     themeID,
		ScheduledAt: scheduledAt,
		Location:    location,
		Notes:       notes,
	}
}
