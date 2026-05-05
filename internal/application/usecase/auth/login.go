package auth

import (
	"context"

	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/mapper"
	"reader-club/internal/domain/repository"
	infraauth "reader-club/internal/infrastructure/auth"
)

type Login struct {
	users  repository.UserRepository
	tokens infraauth.TokenService
}

func NewLogin(users repository.UserRepository, tokens infraauth.TokenService) *Login {
	return &Login{users: users, tokens: tokens}
}

func (uc *Login) Execute(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := uc.users.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperr.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, apperr.ErrInvalidCredentials
	}

	token, err := uc.tokens.Generate(user.ID.String(), user.Email, []string{string(user.GlobalRole)})
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  mapper.ToUserResponse(user),
	}, nil
}
