package usecase

import (
	"context"
	"time"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type MatchService struct {
	repo ports.MatchRepo
}

func NewMatchService(repo ports.MatchRepo) *MatchService {
	return &MatchService{repo: repo}
}

func (s *MatchService) Create(ctx context.Context, sessionID string, difficultyID, level int) (entity.Match, error) {
	match := entity.Match{
		SessionID:    sessionID,
		DifficultyID: difficultyID,
		LevelN:       level,
		IsActive:     true,
	}
	return s.repo.Create(ctx, match)
}

func (s *MatchService) GetActiveBySession(ctx context.Context, sessionID string) (entity.Match, error) {
	return s.repo.GetActiveBySession(ctx, sessionID)
}

func (s *MatchService) Get(ctx context.Context, id string) (entity.Match, error) {
	return s.repo.Get(ctx, id)
}

func (s *MatchService) FinishActive(ctx context.Context, sessionID string, outcome *string) (entity.Match, error) {
	match, err := s.repo.GetActiveBySession(ctx, sessionID)
	if err != nil {
		return entity.Match{}, err
	}
	now := time.Now().UTC()
	match.IsActive = false
	match.EndedAt = &now
	if outcome != nil {
		match.Outcome = outcome
	}
	return s.repo.Update(ctx, match)
}
