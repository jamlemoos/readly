package dto

import (
	"time"

	"github.com/google/uuid"
	"reader-club/internal/domain/entity"
)

type JoinClubRequest struct {
	Role entity.Role `json:"role" binding:"required"`
}

type MembershipResponse struct {
	ID       uuid.UUID    `json:"id"`
	User     UserResponse `json:"user"`
	Role     entity.Role  `json:"role"`
	JoinedAt time.Time    `json:"joined_at"`
}
