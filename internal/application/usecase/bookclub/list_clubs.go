package bookclub

import (
	"context"

	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/repository"
)

type ListBookClubs struct {
	clubs repository.BookClubRepository
}

func NewListBookClubs(clubs repository.BookClubRepository) *ListBookClubs {
	return &ListBookClubs{clubs: clubs}
}

func (uc *ListBookClubs) Execute(ctx context.Context) ([]dto.BookClubSummaryResponse, error) {
	clubs, err := uc.clubs.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BookClubSummaryResponse, len(clubs))
	for i, c := range clubs {
		result[i] = mapper.ToBookClubSummaryResponse(c)
	}
	return result, nil
}
