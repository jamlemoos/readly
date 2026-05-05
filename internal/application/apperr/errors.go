package apperr

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoEligibleBooks    = errors.New("no eligible books to draw from")
	ErrThemeAlreadyExists = errors.New("a theme already exists for this month")
	ErrMemberLimitReached = errors.New("club member limit reached")
)
