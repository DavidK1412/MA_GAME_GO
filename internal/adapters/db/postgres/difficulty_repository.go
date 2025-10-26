package postgres

import (
	"context"

	"github.com/org/ranas-bdi-backend/internal/domain/entity"
	"github.com/org/ranas-bdi-backend/internal/domain/ports"
)

type DifficultyRepository struct {
	pool pgxQuerier
}

var _ ports.DifficultyRepo = (*DifficultyRepository)(nil)

func NewDifficultyRepository(pool pgxQuerier) *DifficultyRepository {
	return &DifficultyRepository{pool: pool}
}

func (r *DifficultyRepository) GetByID(ctx context.Context, id int) (entity.Difficulty, error) {
	var difficulty entity.Difficulty
	query := `
        SELECT id, name, number_of_blocks
        FROM difficulty
        WHERE id = $1
    `
	row := r.pool.QueryRow(ctx, query, id)
	if err := row.Scan(&difficulty.ID, &difficulty.Name, &difficulty.NumberOfBlocks); err != nil {
		return entity.Difficulty{}, err
	}
	return difficulty, nil
}

func (r *DifficultyRepository) GetAll(ctx context.Context) ([]entity.Difficulty, error) {
	query := `
        SELECT id, name, number_of_blocks
        FROM difficulty
        ORDER BY id
    `
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var difficulties []entity.Difficulty
	for rows.Next() {
		var d entity.Difficulty
		if err := rows.Scan(&d.ID, &d.Name, &d.NumberOfBlocks); err != nil {
			return nil, err
		}
		difficulties = append(difficulties, d)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return difficulties, nil
}
