package httpadapter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/org/ranas-bdi-backend/internal/app/usecase"
	"github.com/org/ranas-bdi-backend/internal/domain/entity"
)

type stubSessionRepo struct {
	lastCreateInput entity.Session
	result          entity.Session
	err             error
}

func (s *stubSessionRepo) Create(ctx context.Context, session entity.Session) (entity.Session, error) {
	s.lastCreateInput = session
	if s.err != nil {
		return entity.Session{}, s.err
	}
	return s.result, nil
}

func (s *stubSessionRepo) Get(ctx context.Context, id string) (entity.Session, error) {
	return entity.Session{}, errors.New("not implemented")
}

func (s *stubSessionRepo) Update(ctx context.Context, session entity.Session) (entity.Session, error) {
	return entity.Session{}, errors.New("not implemented")
}

func TestHandleCreateGame_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubSessionRepo{}
	sessionService := usecase.NewSessionService(repo)

	router := gin.New()
	handler := &Handler{router: router, sessions: sessionService}
	router.POST("/game", handler.handleCreateGame)

	gameID := uuid.NewString()
	startedAt := time.Now().UTC().Round(0)
	expectedPlayerID := gameID
	repo.result = entity.Session{
		ID:         "session-123",
		PlayerID:   &expectedPlayerID,
		Device:     nil,
		IsFinished: false,
		StartedAt:  startedAt,
		EndedAt:    nil,
	}

	payload, err := json.Marshal(map[string]string{"game_id": gameID})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/game", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var body entity.Session
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.Equal(t, repo.result, body)

	require.NotNil(t, repo.lastCreateInput.PlayerID)
	require.Equal(t, gameID, *repo.lastCreateInput.PlayerID)
	require.Nil(t, repo.lastCreateInput.Device)
}

func TestHandleCreateGame_InvalidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := &stubSessionRepo{}
	sessionService := usecase.NewSessionService(repo)

	router := gin.New()
	handler := &Handler{router: router, sessions: sessionService}
	router.POST("/game", handler.handleCreateGame)

	payload, err := json.Marshal(map[string]string{"game_id": "not-a-uuid"})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/game", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code)

	var body map[string]string
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	require.Equal(t, errInvalidGameID.Error(), body["error"])
	require.Empty(t, repo.lastCreateInput.ID)
}
