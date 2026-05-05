package bookclub

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

type CreateBookClub struct {
	clubs     repository.BookClubRepository
	members   repository.MembershipRepository
	auditLogs repository.AuditLogRepository
	logger    applogger.Logger
}

func NewCreateBookClub(
	clubs repository.BookClubRepository,
	members repository.MembershipRepository,
	auditLogs repository.AuditLogRepository,
	logger applogger.Logger,
) *CreateBookClub {
	return &CreateBookClub{clubs: clubs, members: members, auditLogs: auditLogs, logger: logger}
}

func (uc *CreateBookClub) Execute(ctx context.Context, actorID uuid.UUID, req dto.CreateBookClubRequest) (*dto.BookClubResponse, error) {
	exists, err := uc.clubs.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.ErrAlreadyExists
	}

	club := entity.NewBookClub(req.Name, req.Description, actorID)
	if err := uc.clubs.Save(ctx, club); err != nil {
		return nil, err
	}

	membership := entity.NewMembership(actorID, club.ID, entity.RoleAdmin)
	if err := uc.members.Save(ctx, membership); err != nil {
		return nil, err
	}

	if al, err := entity.NewAuditLog(actorID, "CREATE_CLUB", "BookClub", club.ID, map[string]string{"name": club.Name}); err == nil {
		if saveErr := uc.auditLogs.Save(ctx, al); saveErr != nil {
			uc.logger.Warn("failed to persist CREATE_CLUB audit log", saveErr)
		}
	}

	club.Owner = entity.User{ID: actorID}
	resp := mapper.ToBookClubResponse(club)
	return &resp, nil
}
