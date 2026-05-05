package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookSuggestion struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	ClubID      uuid.UUID `gorm:"type:uuid;not null"`
	SuggestedBy uuid.UUID `gorm:"type:uuid;not null"`
	Title       string    `gorm:"not null"`
	Author      string    `gorm:"not null"`
	Description string
	SuggestedAt time.Time

	Club BookClub `gorm:"foreignKey:ClubID"`
	User User     `gorm:"foreignKey:SuggestedBy"`
}

func NewBookSuggestion(clubID, userID uuid.UUID, title, author, description string) *BookSuggestion {
	return &BookSuggestion{
		ID:          uuid.New(),
		ClubID:      clubID,
		SuggestedBy: userID,
		Title:       title,
		Author:      author,
		Description: description,
		SuggestedAt: time.Now(),
	}
}
