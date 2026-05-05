package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToMonthlyThemeResponse(t *entity.MonthlyTheme) dto.MonthlyThemeResponse {
	return dto.MonthlyThemeResponse{
		ID:             t.ID,
		Year:           t.Year,
		Month:          t.Month,
		BookSuggestion: ToBookSuggestionResponse(&t.BookSuggestion),
		DrawnAt:        t.DrawnAt,
	}
}
