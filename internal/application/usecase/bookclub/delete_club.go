package bookclub

import (
	"context"

	"reader-club/internal/application/apperr"
	applogger "reader-club/internal/application/logger"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"

	"github.com/google/uuid"
)

type DeleteBookClub struct {
	clubs       repository.BookClubRepository
	meetings    repository.MeetingRepository
	themes      repository.MonthlyThemeRepository
	suggestions repository.BookSuggestionRepository
	members     repository.MembershipRepository
	auditLogs   repository.AuditLogRepository
	logger      applogger.Logger
}

func NewDeleteBookClub(
	clubs repository.BookClubRepository,
	meetings repository.MeetingRepository,
	themes repository.MonthlyThemeRepository,
	suggestions repository.BookSuggestionRepository,
	members repository.MembershipRepository,
	auditLogs repository.AuditLogRepository,
	logger applogger.Logger,
) *DeleteBookClub {
	return &DeleteBookClub{
		clubs:       clubs,
		meetings:    meetings,
		themes:      themes,
		suggestions: suggestions,
		members:     members,
		auditLogs:   auditLogs,
		logger:      logger,
	}
}

func (uc *DeleteBookClub) Execute(ctx context.Context, actorID, clubID uuid.UUID) error {
	if _, err := uc.clubs.FindByID(ctx, clubID); err != nil {
		return apperr.ErrNotFound
	}

	if err := uc.meetings.DeleteByClubID(ctx, clubID); err != nil {
		return err
	}
	if err := uc.themes.DeleteByClubID(ctx, clubID); err != nil {
		return err
	}
	if err := uc.suggestions.DeleteByClubID(ctx, clubID); err != nil {
		return err
	}
	if err := uc.members.DeleteByClubID(ctx, clubID); err != nil {
		return err
	}
	if err := uc.clubs.Delete(ctx, clubID); err != nil {
		return err
	}

	if al, err := entity.NewAuditLog(actorID, "DELETE_CLUB", "BookClub", clubID, nil); err == nil {
		if saveErr := uc.auditLogs.Save(ctx, al); saveErr != nil {
			uc.logger.Warn("failed to persist DELETE_CLUB audit log", saveErr)
		}
	}

	return nil
}
