package usecase

import (
	"context"
	"time"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type SessionService struct {
	repo ports.SessionRepo
}

func NewSessionService(repo ports.SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(ctx context.Context, playerID, device string) (entity.Session, error) {
	session := entity.Session{
		IsFinished: false,
	}
	if playerID != "" {
		session.PlayerID = &playerID
	}
	if device != "" {
		session.Device = &device
	}
	return s.repo.Create(ctx, session)
}

func (s *SessionService) Get(ctx context.Context, id string) (entity.Session, error) {
	return s.repo.Get(ctx, id)
}

func (s *SessionService) Finish(ctx context.Context, id string) (entity.Session, error) {
	session, err := s.repo.Get(ctx, id)
	if err != nil {
		return entity.Session{}, err
	}
	now := time.Now().UTC()
	session.IsFinished = true
	session.EndedAt = &now
	return s.repo.Update(ctx, session)
}
