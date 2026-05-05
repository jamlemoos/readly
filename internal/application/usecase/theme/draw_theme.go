package theme

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	applogger "reader-club/internal/application/logger"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type DrawMonthlyTheme struct {
	members     repository.MembershipRepository
	suggestions repository.BookSuggestionRepository
	themes      repository.MonthlyThemeRepository
	auditLogs   repository.AuditLogRepository
	logger      applogger.Logger
}

func NewDrawMonthlyTheme(
	members repository.MembershipRepository,
	suggestions repository.BookSuggestionRepository,
	themes repository.MonthlyThemeRepository,
	auditLogs repository.AuditLogRepository,
	logger applogger.Logger,
) *DrawMonthlyTheme {
	return &DrawMonthlyTheme{members: members, suggestions: suggestions, themes: themes, auditLogs: auditLogs, logger: logger}
}

func (uc *DrawMonthlyTheme) Execute(ctx context.Context, actorID, clubID uuid.UUID, req dto.DrawThemeRequest) (*dto.MonthlyThemeResponse, error) {
	membership, err := uc.members.FindByUserAndClub(ctx, actorID, clubID)
	if err != nil || !membership.CanManageClub() {
		return nil, apperr.ErrForbidden
	}

	exists, err := uc.themes.ExistsByClubAndMonth(ctx, clubID, req.Year, req.Month)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.ErrThemeAlreadyExists
	}

	eligible, err := uc.suggestions.FindEligibleForDraw(ctx, clubID)
	if err != nil {
		return nil, err
	}
	if len(eligible) == 0 {
		return nil, apperr.ErrNoEligibleBooks
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(eligible), func(i, j int) {
		eligible[i], eligible[j] = eligible[j], eligible[i]
	})
	drawn := eligible[0]

	theme := entity.NewMonthlyTheme(clubID, drawn.ID, req.Year, req.Month)
	theme.BookSuggestion = drawn

	if err := uc.themes.Save(ctx, theme); err != nil {
		return nil, err
	}

	if al, err := entity.NewAuditLog(actorID, "DRAW_THEME", "MonthlyTheme", theme.ID, map[string]any{
		"year":       req.Year,
		"month":      req.Month,
		"book":       drawn.Title,
		"candidates": len(eligible),
		"message":    fmt.Sprintf("Theme drawn: '%s' from %d candidates", drawn.Title, len(eligible)),
	}); err == nil {
		if saveErr := uc.auditLogs.Save(ctx, al); saveErr != nil {
			uc.logger.Warn("failed to persist DRAW_THEME audit log", saveErr)
		}
	}

	resp := mapper.ToMonthlyThemeResponse(theme)
	return &resp, nil
}
