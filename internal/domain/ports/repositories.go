package ports

import (
	"context"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
)

type SessionRepo interface {
	Create(ctx context.Context, session entity.Session) (entity.Session, error)
	Get(ctx context.Context, id string) (entity.Session, error)
	Update(ctx context.Context, session entity.Session) (entity.Session, error)
}

type MatchRepo interface {
	Create(ctx context.Context, match entity.Match) (entity.Match, error)
	Get(ctx context.Context, id string) (entity.Match, error)
	Update(ctx context.Context, match entity.Match) (entity.Match, error)
	GetActiveBySession(ctx context.Context, sessionID string) (entity.Match, error)
}

type MoveRepo interface {
	Create(ctx context.Context, move entity.Move) (entity.Move, error)
	GetByMatch(ctx context.Context, matchID string) ([]entity.Move, error)
	GetLastByMatch(ctx context.Context, matchID string) (entity.Move, error)
}

type DifficultyRepo interface {
	GetByID(ctx context.Context, id int) (entity.Difficulty, error)
	GetAll(ctx context.Context) ([]entity.Difficulty, error)
}
