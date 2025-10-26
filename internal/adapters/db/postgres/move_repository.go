package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type MoveRepository struct {
	pool pgxQuerier
}

var _ ports.MoveRepo = (*MoveRepository)(nil)

func NewMoveRepository(pool pgxQuerier) *MoveRepository {
	return &MoveRepository{pool: pool}
}

func (r *MoveRepository) Create(ctx context.Context, move entity.Move) (entity.Move, error) {
	var created entity.Move
	query := `
        INSERT INTO moves (
            match_id, seq, occurred_at, elapsed_ms, from_idx, to_idx,
            move_kind, frog_side, is_correct, interruption,
            board_before, board_after, branching_factor, buclicidad
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
        RETURNING id, match_id, seq, occurred_at, elapsed_ms, from_idx, to_idx,
                  move_kind, frog_side, is_correct, interruption,
                  board_before, board_after, branching_factor, buclicidad
    `
	row := r.pool.QueryRow(ctx, query,
		move.MatchID,
		move.Seq,
		move.OccurredAt,
		move.ElapsedMs,
		move.FromIdx,
		move.ToIdx,
		move.MoveKind,
		move.FrogSide,
		move.IsCorrect,
		move.Interruption,
		nullableBytes(move.BoardBefore),
		nullableBytes(move.BoardAfter),
		nullableInt(move.BranchingFactor),
		nullableFloat(move.Buclicidad),
	)
	if err := scanMove(row, &created); err != nil {
		return entity.Move{}, err
	}
	return created, nil
}

func (r *MoveRepository) GetByMatch(ctx context.Context, matchID string) ([]entity.Move, error) {
	query := `
        SELECT id, match_id, seq, occurred_at, elapsed_ms, from_idx, to_idx,
               move_kind, frog_side, is_correct, interruption,
               board_before, board_after, branching_factor, buclicidad
        FROM moves
        WHERE match_id = $1
        ORDER BY seq
    `
	rows, err := r.pool.Query(ctx, query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var moves []entity.Move
	for rows.Next() {
		var mv entity.Move
		if err := scanMove(rows, &mv); err != nil {
			return nil, err
		}
		moves = append(moves, mv)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return moves, nil
}

func (r *MoveRepository) GetLastByMatch(ctx context.Context, matchID string) (entity.Move, error) {
	var move entity.Move
	query := `
        SELECT id, match_id, seq, occurred_at, elapsed_ms, from_idx, to_idx,
               move_kind, frog_side, is_correct, interruption,
               board_before, board_after, branching_factor, buclicidad
        FROM moves
        WHERE match_id = $1
        ORDER BY seq DESC
        LIMIT 1
    `
	row := r.pool.QueryRow(ctx, query, matchID)
	if err := scanMove(row, &move); err != nil {
		return entity.Move{}, err
	}
	return move, nil
}

func scanMove(row pgx.Row, move *entity.Move) error {
	var (
		boardBefore []byte
		boardAfter  []byte
		branching   sql.NullInt64
		buclicidad  sql.NullFloat64
	)
	if err := row.Scan(
		&move.ID,
		&move.MatchID,
		&move.Seq,
		&move.OccurredAt,
		&move.ElapsedMs,
		&move.FromIdx,
		&move.ToIdx,
		&move.MoveKind,
		&move.FrogSide,
		&move.IsCorrect,
		&move.Interruption,
		&boardBefore,
		&boardAfter,
		&branching,
		&buclicidad,
	); err != nil {
		return err
	}
	if len(boardBefore) > 0 {
		move.BoardBefore = make([]byte, len(boardBefore))
		copy(move.BoardBefore, boardBefore)
	} else {
		move.BoardBefore = nil
	}
	if len(boardAfter) > 0 {
		move.BoardAfter = make([]byte, len(boardAfter))
		copy(move.BoardAfter, boardAfter)
	} else {
		move.BoardAfter = nil
	}
	move.BranchingFactor = intPtrFromNull(branching)
	move.Buclicidad = floatPtrFromNull(buclicidad)
	return nil
}
