package usecase

import (
	"context"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type MoveService struct {
	repo ports.MoveRepo
}

func NewMoveService(repo ports.MoveRepo) *MoveService {
	return &MoveService{repo: repo}
}

func (s *MoveService) Create(ctx context.Context, move entity.Move) (entity.Move, error) {
	return s.repo.Create(ctx, move)
}

func (s *MoveService) ListByMatch(ctx context.Context, matchID string) ([]entity.Move, error) {
	return s.repo.GetByMatch(ctx, matchID)
}

func (s *MoveService) GetLastByMatch(ctx context.Context, matchID string) (entity.Move, error) {
	return s.repo.GetLastByMatch(ctx, matchID)
}
