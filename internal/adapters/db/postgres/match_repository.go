package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type MatchRepository struct {
	pool pgxQuerier
}

var _ ports.MatchRepo = (*MatchRepository)(nil)

func NewMatchRepository(pool pgxQuerier) *MatchRepository {
	return &MatchRepository{pool: pool}
}

func (r *MatchRepository) Create(ctx context.Context, match entity.Match) (entity.Match, error) {
	var created entity.Match
	query := `
        INSERT INTO matches (session_id, difficulty_id, level_n, is_active, outcome, meta)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, session_id, difficulty_id, level_n, is_active, started_at, ended_at, outcome, meta
    `
	row := r.pool.QueryRow(ctx, query,
		match.SessionID,
		match.DifficultyID,
		match.LevelN,
		match.IsActive,
		nullableString(match.Outcome),
		nullableBytes(match.Meta),
	)
	if err := scanMatch(row, &created); err != nil {
		return entity.Match{}, err
	}
	return created, nil
}

func (r *MatchRepository) Get(ctx context.Context, id string) (entity.Match, error) {
	var match entity.Match
	query := `
        SELECT id, session_id, difficulty_id, level_n, is_active, started_at, ended_at, outcome, meta
        FROM matches
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)
	if err := scanMatch(row, &match); err != nil {
		return entity.Match{}, err
	}
	return match, nil
}

func (r *MatchRepository) Update(ctx context.Context, match entity.Match) (entity.Match, error) {
	var updated entity.Match
	query := `
        UPDATE matches
        SET session_id = $2,
            difficulty_id = $3,
            level_n = $4,
            is_active = $5,
            started_at = $6,
            ended_at = $7,
            outcome = $8,
            meta = $9
        WHERE id = $1
        RETURNING id, session_id, difficulty_id, level_n, is_active, started_at, ended_at, outcome, meta
    `
	row := r.pool.QueryRow(ctx, query,
		match.ID,
		match.SessionID,
		match.DifficultyID,
		match.LevelN,
		match.IsActive,
		match.StartedAt,
		nullableTime(match.EndedAt),
		nullableString(match.Outcome),
		nullableBytes(match.Meta),
	)
	if err := scanMatch(row, &updated); err != nil {
		return entity.Match{}, err
	}
	return updated, nil
}

func (r *MatchRepository) GetActiveBySession(ctx context.Context, sessionID string) (entity.Match, error) {
	var match entity.Match
	query := `
        SELECT id, session_id, difficulty_id, level_n, is_active, started_at, ended_at, outcome, meta
        FROM matches
        WHERE session_id = $1 AND is_active = TRUE
        LIMIT 1
    `
	row := r.pool.QueryRow(ctx, query, sessionID)
	if err := scanMatch(row, &match); err != nil {
		return entity.Match{}, err
	}
	return match, nil
}

func scanMatch(row pgx.Row, match *entity.Match) error {
	var (
		endedAt sql.NullTime
		outcome sql.NullString
		meta    []byte
	)
	if err := row.Scan(
		&match.ID,
		&match.SessionID,
		&match.DifficultyID,
		&match.LevelN,
		&match.IsActive,
		&match.StartedAt,
		&endedAt,
		&outcome,
		&meta,
	); err != nil {
		return err
	}
	match.EndedAt = timePtrFromNull(endedAt)
	match.Outcome = stringPtrFromNull(outcome)
	if len(meta) > 0 {
		match.Meta = make([]byte, len(meta))
		copy(match.Meta, meta)
	} else {
		match.Meta = nil
	}
	return nil
}
