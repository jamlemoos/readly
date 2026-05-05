package dto

import (
	"time"

	"github.com/google/uuid"
)

type DrawThemeRequest struct {
	Year  int `json:"year"  binding:"required,min=2000"`
	Month int `json:"month" binding:"required,min=1,max=12"`
}

type MonthlyThemeResponse struct {
	ID             uuid.UUID              `json:"id"`
	Year           int                    `json:"year"`
	Month          int                    `json:"month"`
	BookSuggestion BookSuggestionResponse `json:"book_suggestion"`
	DrawnAt        time.Time              `json:"drawn_at"`
}
