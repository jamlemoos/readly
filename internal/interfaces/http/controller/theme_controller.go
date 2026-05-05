package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	themeuc "reader-club/internal/application/usecase/theme"
	"reader-club/internal/interfaces/http/httperr"
	"reader-club/internal/interfaces/http/middleware"
)

type ThemeController struct {
	draw *themeuc.DrawMonthlyTheme
}

func NewThemeController(draw *themeuc.DrawMonthlyTheme) *ThemeController {
	return &ThemeController{draw: draw}
}

func (h *ThemeController) Draw(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	var req dto.DrawThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	actorID := middleware.ActorID(c)
	theme, err := h.draw.Execute(c.Request.Context(), actorID, clubID, req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			httperr.Respond(c, http.StatusForbidden, "forbidden", "only admins can draw themes")
		case errors.Is(err, apperr.ErrThemeAlreadyExists):
			httperr.Respond(c, http.StatusConflict, "conflict", "a theme already exists for this month")
		case errors.Is(err, apperr.ErrNoEligibleBooks):
			httperr.Respond(c, http.StatusUnprocessableEntity, "unprocessable_entity", "no eligible book suggestions to draw from")
		default:
			httperr.RespondInternalError(c)
		}
		return
	}

	c.JSON(http.StatusCreated, theme)
}
