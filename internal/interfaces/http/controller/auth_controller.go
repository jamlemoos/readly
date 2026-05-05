package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	authuc "reader-club/internal/application/usecase/auth"
	"reader-club/internal/interfaces/http/httperr"
)

type AuthController struct {
	register *authuc.Register
	login    *authuc.Login
}

func NewAuthController(register *authuc.Register, login *authuc.Login) *AuthController {
	return &AuthController{register: register, login: login}
}

func (h *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	user, err := h.register.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrAlreadyExists) {
			httperr.Respond(c, http.StatusConflict, "conflict", "email already registered")
			return
		}
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	resp, err := h.login.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, apperr.ErrInvalidCredentials) {
			httperr.Respond(c, http.StatusUnauthorized, "unauthorized", "invalid credentials")
			return
		}
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, resp)
}
