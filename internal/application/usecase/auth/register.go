package auth

import (
	"context"

	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/entity"
	"reader-club/internal/domain/repository"
)

type Register struct {
	users repository.UserRepository
}

func NewRegister(users repository.UserRepository) *Register {
	return &Register{users: users}
}

func (uc *Register) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	exists, err := uc.users.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.ErrAlreadyExists
	}

	user, err := entity.NewUser(req.Name, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	if err := uc.users.Save(ctx, user); err != nil {
		return nil, err
	}

	resp := mapper.ToUserResponse(user)
	return &resp, nil
}
