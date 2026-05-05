package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToBookClubResponse(c *entity.BookClub) dto.BookClubResponse {
	return dto.BookClubResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Owner:       ToUserResponse(&c.Owner),
		CreatedAt:   c.CreatedAt,
	}
}

func ToBookClubSummaryResponse(c entity.BookClub) dto.BookClubSummaryResponse {
	return dto.BookClubSummaryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
	}
}
