package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateBookClubRequest struct {
	Name        string `json:"name"        binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
}

type BookClubResponse struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Owner       UserResponse `json:"owner"`
	CreatedAt   time.Time    `json:"created_at"`
}

type BookClubSummaryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
