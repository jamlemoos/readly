package book

import (
	"context"

	"github.com/google/uuid"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/repository"
)

type ListBookSuggestions struct {
	suggestions repository.BookSuggestionRepository
}

func NewListBookSuggestions(suggestions repository.BookSuggestionRepository) *ListBookSuggestions {
	return &ListBookSuggestions{suggestions: suggestions}
}

func (uc *ListBookSuggestions) Execute(ctx context.Context, clubID uuid.UUID) ([]dto.BookSuggestionResponse, error) {
	items, err := uc.suggestions.FindByClub(ctx, clubID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BookSuggestionResponse, len(items))
	for i := range items {
		result[i] = mapper.ToBookSuggestionResponse(&items[i])
	}
	return result, nil
}
