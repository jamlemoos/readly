package meeting

import (
	"context"

	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	applogger "reader-club/internal/application/logger"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type CreateMeeting struct {
	members   repository.MembershipRepository
	themes    repository.MonthlyThemeRepository
	meetings  repository.MeetingRepository
	auditLogs repository.AuditLogRepository
	logger    applogger.Logger
}

func NewCreateMeeting(
	members repository.MembershipRepository,
	themes repository.MonthlyThemeRepository,
	meetings repository.MeetingRepository,
	auditLogs repository.AuditLogRepository,
	logger applogger.Logger,
) *CreateMeeting {
	return &CreateMeeting{members: members, themes: themes, meetings: meetings, auditLogs: auditLogs, logger: logger}
}

func (uc *CreateMeeting) Execute(ctx context.Context, actorID, clubID uuid.UUID, req dto.CreateMeetingRequest) (*dto.MeetingResponse, error) {
	membership, err := uc.members.FindByUserAndClub(ctx, actorID, clubID)
	if err != nil || !membership.CanManageClub() {
		return nil, apperr.ErrForbidden
	}

	theme, err := uc.themes.FindByID(ctx, req.ThemeID)
	if err != nil {
		return nil, apperr.ErrNotFound
	}

	meeting := entity.NewMeeting(clubID, req.ThemeID, req.ScheduledAt, req.Location, req.Notes)
	meeting.Theme = *theme

	if err := uc.meetings.Save(ctx, meeting); err != nil {
		return nil, err
	}

	if al, err := entity.NewAuditLog(actorID, "CREATE_MEETING", "Meeting", meeting.ID, map[string]any{
		"scheduled_at": req.ScheduledAt,
	}); err == nil {
		if saveErr := uc.auditLogs.Save(ctx, al); saveErr != nil {
			uc.logger.Warn("failed to persist CREATE_MEETING audit log", saveErr)
		}
	}

	resp := mapper.ToMeetingResponse(meeting)
	return &resp, nil
}
