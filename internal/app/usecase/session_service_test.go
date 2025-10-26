package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
)

type stubSessionRepo struct {
	createFn func(ctx context.Context, session entity.Session) (entity.Session, error)
	getFn    func(ctx context.Context, id string) (entity.Session, error)
	updateFn func(ctx context.Context, session entity.Session) (entity.Session, error)
}

func (s stubSessionRepo) Create(ctx context.Context, session entity.Session) (entity.Session, error) {
	if s.createFn != nil {
		return s.createFn(ctx, session)
	}
	return entity.Session{}, nil
}

func (s stubSessionRepo) Get(ctx context.Context, id string) (entity.Session, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return entity.Session{}, nil
}

func (s stubSessionRepo) Update(ctx context.Context, session entity.Session) (entity.Session, error) {
	if s.updateFn != nil {
		return s.updateFn(ctx, session)
	}
	return entity.Session{}, nil
}

func TestSessionServiceCreate(t *testing.T) {
	ctx := context.Background()
	captured := entity.Session{}
	repo := stubSessionRepo{
		createFn: func(ctx context.Context, session entity.Session) (entity.Session, error) {
			captured = session
			session.ID = "session-id"
			session.StartedAt = time.Now().UTC()
			return session, nil
		},
	}

	svc := NewSessionService(repo)
	session, err := svc.Create(ctx, "player-id", "Meta Quest 3")
	require.NoError(t, err)
	require.Equal(t, "session-id", session.ID)
	require.False(t, session.IsFinished)
	require.NotNil(t, captured.PlayerID)
	require.Equal(t, "player-id", *captured.PlayerID)
	require.NotNil(t, captured.Device)
	require.Equal(t, "Meta Quest 3", *captured.Device)
}
