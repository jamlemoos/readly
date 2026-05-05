package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMeetingRequest struct {
	ThemeID     uuid.UUID `json:"theme_id"     binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	Location    string    `json:"location"     binding:"max=200"`
	Notes       string    `json:"notes"        binding:"max=2000"`
}

type MeetingResponse struct {
	ID          uuid.UUID            `json:"id"`
	Theme       MonthlyThemeResponse `json:"theme"`
	ScheduledAt time.Time            `json:"scheduled_at"`
	Location    string               `json:"location"`
	Notes       string               `json:"notes"`
	CreatedAt   time.Time            `json:"created_at"`
}
