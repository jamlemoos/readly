package mapper

import (
	"reader-club/internal/application/dto"
	"reader-club/internal/domain/entity"
)

func ToMeetingResponse(m *entity.Meeting) dto.MeetingResponse {
	return dto.MeetingResponse{
		ID:          m.ID,
		Theme:       ToMonthlyThemeResponse(&m.Theme),
		ScheduledAt: m.ScheduledAt,
		Location:    m.Location,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
	}
}
