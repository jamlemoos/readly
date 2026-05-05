package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	membershipuc "reader-club/internal/application/usecase/membership"
	"reader-club/internal/interfaces/http/httperr"
	"reader-club/internal/interfaces/http/middleware"
)

type MembershipController struct {
	join *membershipuc.JoinClub
}

func NewMembershipController(join *membershipuc.JoinClub) *MembershipController {
	return &MembershipController{join: join}
}

func (h *MembershipController) Join(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	var req dto.JoinClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	actorID := middleware.ActorID(c)
	membership, err := h.join.Execute(c.Request.Context(), actorID, clubID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrAlreadyExists):
			httperr.Respond(c, http.StatusConflict, "conflict", "already a member of this club")
		case errors.Is(err, apperr.ErrNotFound):
			httperr.Respond(c, http.StatusNotFound, "not_found", "club not found")
		case errors.Is(err, apperr.ErrMemberLimitReached):
			httperr.Respond(c, http.StatusUnprocessableEntity, "unprocessable_entity", "club member limit has been reached")
		default:
			httperr.RespondInternalError(c)
		}
		return
	}

	c.JSON(http.StatusCreated, membership)
}
