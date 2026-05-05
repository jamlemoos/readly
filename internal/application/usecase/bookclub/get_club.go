package bookclub

import (
	"context"

	"github.com/google/uuid"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/repository"
)

type GetBookClub struct {
	clubs repository.BookClubRepository
}

func NewGetBookClub(clubs repository.BookClubRepository) *GetBookClub {
	return &GetBookClub{clubs: clubs}
}

func (uc *GetBookClub) Execute(ctx context.Context, id uuid.UUID) (*dto.BookClubResponse, error) {
	club, err := uc.clubs.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := mapper.ToBookClubResponse(club)
	return &resp, nil
}
