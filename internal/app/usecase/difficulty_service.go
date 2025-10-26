package usecase

import (
	"context"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type DifficultyService struct {
	repo ports.DifficultyRepo
}

func NewDifficultyService(repo ports.DifficultyRepo) *DifficultyService {
	return &DifficultyService{repo: repo}
}

func (s *DifficultyService) GetAll(ctx context.Context) ([]entity.Difficulty, error) {
	return s.repo.GetAll(ctx)
}

func (s *DifficultyService) GetByID(ctx context.Context, id int) (entity.Difficulty, error) {
	return s.repo.GetByID(ctx, id)
}
