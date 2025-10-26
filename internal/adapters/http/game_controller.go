package httpadapter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createGameReq struct {
	GameID string `json:"game_id" binding:"required"`
}

var errInvalidGameID = errors.New("game_id must be a valid UUID")

func (h *Handler) handleCreateGame(c *gin.Context) {
	var req createGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if _, err := uuid.Parse(req.GameID); err != nil {
		respondError(c, http.StatusBadRequest, errInvalidGameID)
		return
	}

	session, err := h.sessions.Create(c.Request.Context(), req.GameID, "")
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, session)
}
