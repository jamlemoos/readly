package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	"reader-club/internal/application/usecase/bookclub"
	"reader-club/internal/interfaces/http/httperr"
	"reader-club/internal/interfaces/http/middleware"
)

type BookClubController struct {
	create *bookclub.CreateBookClub
	get    *bookclub.GetBookClub
	list   *bookclub.ListBookClubs
	delete *bookclub.DeleteBookClub
}

func NewBookClubController(
	create *bookclub.CreateBookClub,
	get *bookclub.GetBookClub,
	list *bookclub.ListBookClubs,
	del *bookclub.DeleteBookClub,
) *BookClubController {
	return &BookClubController{create: create, get: get, list: list, delete: del}
}

func (h *BookClubController) Create(c *gin.Context) {
	var req dto.CreateBookClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	actorID := middleware.ActorID(c)
	club, err := h.create.Execute(c.Request.Context(), actorID, req)
	if err != nil {
		if errors.Is(err, apperr.ErrAlreadyExists) {
			httperr.Respond(c, http.StatusConflict, "conflict", "a club with this name already exists")
			return
		}
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusCreated, club)
}

func (h *BookClubController) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	club, err := h.get.Execute(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperr.ErrNotFound) {
			httperr.Respond(c, http.StatusNotFound, "not_found", "club not found")
			return
		}
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, club)
}

func (h *BookClubController) List(c *gin.Context) {
	clubs, err := h.list.Execute(c.Request.Context())
	if err != nil {
		httperr.RespondInternalError(c)
		return
	}
	c.JSON(http.StatusOK, clubs)
}

func (h *BookClubController) Delete(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	actorID := middleware.ActorID(c)
	if err := h.delete.Execute(c.Request.Context(), actorID, clubID); err != nil {
		switch {
		case errors.Is(err, apperr.ErrForbidden):
			httperr.Respond(c, http.StatusForbidden, "forbidden", "only admins can delete clubs")
		case errors.Is(err, apperr.ErrNotFound):
			httperr.Respond(c, http.StatusNotFound, "not_found", "club not found")
		default:
			httperr.RespondInternalError(c)
		}
		return
	}

	c.Status(http.StatusNoContent)
}
