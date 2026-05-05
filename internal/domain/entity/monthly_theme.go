package entity

import (
	"time"

	"github.com/google/uuid"
)

type MonthlyTheme struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClubID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_club_month"`
	BookSuggestionID uuid.UUID `gorm:"type:uuid;not null"`
	Year             int       `gorm:"not null;uniqueIndex:idx_club_month"`
	Month            int       `gorm:"not null;uniqueIndex:idx_club_month"`
	DrawnAt          time.Time

	Club           BookClub       `gorm:"foreignKey:ClubID"`
	BookSuggestion BookSuggestion `gorm:"foreignKey:BookSuggestionID"`
}

func NewMonthlyTheme(clubID, bookSuggestionID uuid.UUID, year, month int) *MonthlyTheme {
	return &MonthlyTheme{
		ID:               uuid.New(),
		ClubID:           clubID,
		BookSuggestionID: bookSuggestionID,
		Year:             year,
		Month:            month,
		DrawnAt:          time.Now(),
	}
}
