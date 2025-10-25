package db

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	AppName         string
}

var (
	pool    *pgxpool.Pool
	initErr error
	once    sync.Once
)

func Init(ctx context.Context, cfg Config) error {
	once.Do(func() {
		if cfg.DSN == "" {
			initErr = errors.New("db.Init: empty DSN")
			return
		}
		pc, err := pgxpool.ParseConfig(cfg.DSN)
		if err != nil {
			initErr = err
			return
		}
		if cfg.AppName != "" {
			pc.ConnConfig.RuntimeParams["application_name"] = cfg.AppName
		}
		if cfg.MaxConns > 0 {
			pc.MaxConns = cfg.MaxConns
		}
		if cfg.MinConns >= 0 {
			pc.MinConns = cfg.MinConns
		}
		if cfg.MaxConnLifetime > 0 {
			pc.MaxConnLifetime = cfg.MaxConnLifetime
		}
		if cfg.MaxConnIdleTime > 0 {
			pc.MaxConnIdleTime = cfg.MaxConnIdleTime
		}

		pool, initErr = pgxpool.NewWithConfig(ctx, pc)
		if initErr == nil {
			ctx2, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			initErr = pool.Ping(ctx2)
			if initErr != nil {
				pool.Close()
				pool = nil
			}
		}
	})
	return initErr
}

func Get() (*pgxpool.Pool, error) {
	if initErr != nil {
		return nil, initErr
	}
	if pool == nil {
		return nil, errors.New("db.Get: pool not initialized; call db.Init first")
	}
	return pool, nil
}

func MustGet() *pgxpool.Pool {
	p, err := Get()
	if err != nil {
		panic(err)
	}
	return p
}

func Close() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}

func WithTx(ctx context.Context, opts pgx.TxOptions, fn func(tx pgx.Tx) error) error {
	p := MustGet()
	tx, err := p.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }() // no-op si ya se hizo Commit

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

type (
	CommandTag = pgconn.CommandTag
	Row        = pgx.Row
	Rows       = pgx.Rows
)
