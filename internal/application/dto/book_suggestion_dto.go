package dto

import (
	"time"

	"github.com/google/uuid"
)

type SuggestBookRequest struct {
	Title       string `json:"title"       binding:"required,min=1,max=200"`
	Author      string `json:"author"      binding:"required,min=1,max=200"`
	Description string `json:"description" binding:"max=1000"`
}

type BookSuggestionResponse struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Author      string       `json:"author"`
	Description string       `json:"description"`
	SuggestedBy UserResponse `json:"suggested_by"`
	SuggestedAt time.Time    `json:"suggested_at"`
}
