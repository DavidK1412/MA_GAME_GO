package httpadapter

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/org/ranas-bdi-backend/internal/app/usecase"
	"github.com/org/ranas-bdi-backend/internal/domain/entity"
)

type Handler struct {
	router        *gin.Engine
	sessions      *usecase.SessionService
	matches       *usecase.MatchService
	moves         *usecase.MoveService
	difficulties  *usecase.DifficultyService
	defaultDevice string
	defaultLevel  int
}

func NewHandler(
	sessions *usecase.SessionService,
	matches *usecase.MatchService,
	moves *usecase.MoveService,
	difficulties *usecase.DifficultyService,
) *Handler {
	router := gin.Default()

	h := &Handler{
		router:        router,
		sessions:      sessions,
		matches:       matches,
		moves:         moves,
		difficulties:  difficulties,
		defaultDevice: "Meta Quest 3",
		defaultLevel:  1,
	}

	h.registerRoutes()
	return h
}

func (h *Handler) registerRoutes() {
	h.router.POST("/sessions", h.handleCreateSession)
	h.router.POST("/matches", h.handleCreateMatch)
	h.router.POST("/matches/:matchID/moves", h.handleCreateMove)
}

func (h *Handler) Router() *gin.Engine {
	return h.router
}

type createSessionRequest struct {
	GameID string `json:"game_id" binding:"required"`
}

type createMatchRequest struct {
	SessionID    string `json:"session_id" binding:"required"`
	DifficultyID int    `json:"difficulty_id" binding:"required"`
}

type createMoveRequest struct {
	Movement []int `json:"movement"`
}

func (h *Handler) handleCreateSession(c *gin.Context) {
	var req createSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	session, err := h.sessions.Create(c.Request.Context(), req.GameID, h.defaultDevice)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *Handler) handleCreateMatch(c *gin.Context) {
	var req createMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if _, err := h.sessions.Get(c.Request.Context(), req.SessionID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(c, http.StatusNotFound, err)
		} else {
			respondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	if _, err := h.difficulties.GetByID(c.Request.Context(), req.DifficultyID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(c, http.StatusNotFound, err)
		} else {
			respondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	match, err := h.matches.Create(c.Request.Context(), req.SessionID, req.DifficultyID, h.defaultLevel)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, match)
}

func (h *Handler) handleCreateMove(c *gin.Context) {
	matchID := c.Param("matchID")
	if matchID == "" {
		respondError(c, http.StatusBadRequest, errMissingMatchID)
		return
	}

	var req createMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	match, err := h.matches.Get(c.Request.Context(), matchID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(c, http.StatusNotFound, err)
		} else {
			respondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	seq := 1
	lastMove, err := h.moves.GetLastByMatch(c.Request.Context(), matchID)
	if err == nil {
		seq = lastMove.Seq + 1
	} else if !errors.Is(err, pgx.ErrNoRows) {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	move := entity.Move{
		MatchID:      match.ID,
		Seq:          seq,
		OccurredAt:   time.Now().UTC(),
		ElapsedMs:    0,
		FromIdx:      0,
		ToIdx:        0,
		MoveKind:     1,
		FrogSide:     1,
		IsCorrect:    true,
		Interruption: false,
	}

	created, err := h.moves.Create(c.Request.Context(), move)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

var errMissingMatchID = errors.New("match id is required")

func respondError(c *gin.Context, status int, err error) {
	if err == nil {
		err = errors.New("unknown error")
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
