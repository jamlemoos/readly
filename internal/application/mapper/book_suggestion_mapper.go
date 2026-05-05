package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToBookSuggestionResponse(s *entity.BookSuggestion) dto.BookSuggestionResponse {
	return dto.BookSuggestionResponse{
		ID:          s.ID,
		Title:       s.Title,
		Author:      s.Author,
		Description: s.Description,
		SuggestedBy: ToUserResponse(&s.User),
		SuggestedAt: s.SuggestedAt,
	}
}
