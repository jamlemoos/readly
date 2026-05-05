package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	meetinguc "reader-club/internal/application/usecase/meeting"
	"reader-club/internal/interfaces/http/httperr"
	"reader-club/internal/interfaces/http/middleware"
)

type MeetingController struct {
	create *meetinguc.CreateMeeting
}

func NewMeetingController(create *meetinguc.CreateMeeting) *MeetingController {
	return &MeetingController{create: create}
}

func (h *MeetingController) Create(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	var req dto.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	actorID := middleware.ActorID(c)
	meeting, err := h.create.Execute(c.Request.Context(), actorID, clubID, req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			httperr.Respond(c, http.StatusForbidden, "forbidden", "only admins can create meetings")
		case errors.Is(err, apperr.ErrNotFound):
			httperr.Respond(c, http.StatusNotFound, "not_found", "theme not found")
		default:
			httperr.RespondInternalError(c)
		}
		return
	}

	c.JSON(http.StatusCreated, meeting)
}
