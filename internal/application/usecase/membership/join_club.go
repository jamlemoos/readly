package membership

import (
	"context"

	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type JoinClub struct {
	clubs   repository.BookClubRepository
	members repository.MembershipRepository
	users   repository.UserRepository
}

func NewJoinClub(
	clubs repository.BookClubRepository,
	members repository.MembershipRepository,
	users repository.UserRepository,
) *JoinClub {
	return &JoinClub{clubs: clubs, members: members, users: users}
}

func (uc *JoinClub) Execute(ctx context.Context, actorID, clubID uuid.UUID, role entity.Role) (*dto.MembershipResponse, error) {
	if !role.IsValid() || role == entity.RoleAdmin {

		role = entity.RoleVisitor
	}

	club, err := uc.clubs.FindByID(ctx, clubID)
	if err != nil {
		return nil, apperr.ErrNotFound
	}

	exists, err := uc.members.ExistsByUserAndClub(ctx, actorID, clubID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.ErrAlreadyExists
	}

	count, err := uc.members.CountByClub(ctx, clubID)
	if err != nil {
		return nil, err
	}
	if !club.CanAddMember(int(count)) {
		return nil, apperr.ErrMemberLimitReached
	}

	membership := entity.NewMembership(actorID, clubID, role)
	if err := uc.members.Save(ctx, membership); err != nil {
		return nil, err
	}

	user, err := uc.users.FindByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	membership.User = *user

	resp := mapper.ToMembershipResponse(membership)
	return &resp, nil
}
