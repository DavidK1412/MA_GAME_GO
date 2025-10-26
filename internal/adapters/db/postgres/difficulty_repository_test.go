package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

type stubRow struct {
	scanFn func(dest ...any) error
}

func (r stubRow) Scan(dest ...any) error {
	return r.scanFn(dest...)
}

type stubQuerier struct {
	row stubRow
}

func (s stubQuerier) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, errors.New("not implemented")
}

func (s stubQuerier) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return s.row
}

func TestDifficultyRepositoryGetByID(t *testing.T) {
	row := stubRow{scanFn: func(dest ...any) error {
		*(dest[0].(*int)) = 2
		*(dest[1].(*string)) = "medium"
		*(dest[2].(*int)) = 9
		return nil
	}}
	repo := NewDifficultyRepository(stubQuerier{row: row})

	diff, err := repo.GetByID(context.Background(), 2)
	require.NoError(t, err)
	require.Equal(t, 2, diff.ID)
	require.Equal(t, "medium", diff.Name)
	require.Equal(t, 9, diff.NumberOfBlocks)
}
