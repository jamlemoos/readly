package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToUserResponse(u *entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		GlobalRole: string(u.GlobalRole),
		CreatedAt:  u.CreatedAt,
	}
}
