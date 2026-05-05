package book

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

type SuggestBook struct {
	members     repository.MembershipRepository
	suggestions repository.BookSuggestionRepository
	users       repository.UserRepository
	auditLogs   repository.AuditLogRepository
	logger      applogger.Logger
}

func NewSuggestBook(
	members repository.MembershipRepository,
	suggestions repository.BookSuggestionRepository,
	users repository.UserRepository,
	auditLogs repository.AuditLogRepository,
	logger applogger.Logger,
) *SuggestBook {
	return &SuggestBook{members: members, suggestions: suggestions, users: users, auditLogs: auditLogs, logger: logger}
}

func (uc *SuggestBook) Execute(ctx context.Context, actorID, clubID uuid.UUID, req dto.SuggestBookRequest) (*dto.BookSuggestionResponse, error) {
	membership, err := uc.members.FindByUserAndClub(ctx, actorID, clubID)
	if err != nil || !membership.CanSuggestBook() {
		return nil, apperr.ErrForbidden
	}

	suggestion := entity.NewBookSuggestion(clubID, actorID, req.Title, req.Author, req.Description)
	if err := uc.suggestions.Save(ctx, suggestion); err != nil {
		return nil, err
	}

	user, err := uc.users.FindByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	suggestion.User = *user

	if al, err := entity.NewAuditLog(actorID, "SUGGEST_BOOK", "BookSuggestion", suggestion.ID, map[string]string{"title": suggestion.Title}); err == nil {
		if saveErr := uc.auditLogs.Save(ctx, al); saveErr != nil {
			uc.logger.Warn("failed to persist SUGGEST_BOOK audit log", saveErr)
		}
	}

	resp := mapper.ToBookSuggestionResponse(suggestion)
	return &resp, nil
}
