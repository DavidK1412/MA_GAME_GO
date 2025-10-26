package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type SessionRepository struct {
	pool pgxQuerier
}

var _ ports.SessionRepo = (*SessionRepository)(nil)

func NewSessionRepository(pool pgxQuerier) *SessionRepository {
	return &SessionRepository{pool: pool}
}

func (r *SessionRepository) Create(ctx context.Context, session entity.Session) (entity.Session, error) {
	var created entity.Session
	query := `
        INSERT INTO sessions (player_id, device, is_finished)
        VALUES ($1, $2, $3)
        RETURNING id, player_id, device, is_finished, started_at, ended_at
    `
	row := r.pool.QueryRow(ctx, query, nullableString(session.PlayerID), nullableString(session.Device), session.IsFinished)
	if err := scanSession(row, &created); err != nil {
		return entity.Session{}, err
	}
	return created, nil
}

func (r *SessionRepository) Get(ctx context.Context, id string) (entity.Session, error) {
	var session entity.Session
	query := `
        SELECT id, player_id, device, is_finished, started_at, ended_at
        FROM sessions
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)
	if err := scanSession(row, &session); err != nil {
		return entity.Session{}, err
	}
	return session, nil
}

func (r *SessionRepository) Update(ctx context.Context, session entity.Session) (entity.Session, error) {
	var updated entity.Session
	query := `
        UPDATE sessions
        SET player_id = $2,
            device = $3,
            is_finished = $4,
            ended_at = $5
        WHERE id = $1
        RETURNING id, player_id, device, is_finished, started_at, ended_at
    `
	row := r.pool.QueryRow(ctx, query,
		session.ID,
		nullableString(session.PlayerID),
		nullableString(session.Device),
		session.IsFinished,
		nullableTime(session.EndedAt),
	)
	if err := scanSession(row, &updated); err != nil {
		return entity.Session{}, err
	}
	return updated, nil
}

func scanSession(row pgx.Row, session *entity.Session) error {
	var (
		playerID sql.NullString
		device   sql.NullString
		endedAt  sql.NullTime
	)
	if err := row.Scan(
		&session.ID,
		&playerID,
		&device,
		&session.IsFinished,
		&session.StartedAt,
		&endedAt,
	); err != nil {
		return err
	}
	session.PlayerID = stringPtrFromNull(playerID)
	session.Device = stringPtrFromNull(device)
	session.EndedAt = timePtrFromNull(endedAt)
	return nil
}
