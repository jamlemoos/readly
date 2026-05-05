package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"reader-club/internal/application/apperr"
	"reader-club/internal/application/dto"
	bookuc "reader-club/internal/application/usecase/book"
	"reader-club/internal/interfaces/http/httperr"
	"reader-club/internal/interfaces/http/middleware"
)

type BookSuggestionController struct {
	suggest *bookuc.SuggestBook
	list    *bookuc.ListBookSuggestions
}

func NewBookSuggestionController(suggest *bookuc.SuggestBook, list *bookuc.ListBookSuggestions) *BookSuggestionController {
	return &BookSuggestionController{suggest: suggest, list: list}
}

func (h *BookSuggestionController) Suggest(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	var req dto.SuggestBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httperr.RespondValidation(c, err)
		return
	}

	actorID := middleware.ActorID(c)
	suggestion, err := h.suggest.Execute(c.Request.Context(), actorID, clubID, req)
	if err != nil {
		if errors.Is(err, apperr.ErrForbidden) {
			httperr.Respond(c, http.StatusForbidden, "forbidden", "only members or admins can suggest books")
			return
		}
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusCreated, suggestion)
}

func (h *BookSuggestionController) List(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("club_id"))
	if err != nil {
		httperr.Respond(c, http.StatusBadRequest, "validation_error", "invalid club_id")
		return
	}

	suggestions, err := h.list.Execute(c.Request.Context(), clubID)
	if err != nil {
		httperr.RespondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, suggestions)
}
