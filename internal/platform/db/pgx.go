package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func InitFromEnv(ctx context.Context) error {
	cfg := Config{
		DSN:     os.Getenv("DATABASE_URL"),
		AppName: os.Getenv("APP_NAME"),
	}

	if max := os.Getenv("DB_MAX_CONNS"); max != "" {
		if v, err := strconv.Atoi(max); err == nil {
			cfg.MaxConns = int32(v)
		}
	}
	if min := os.Getenv("DB_MIN_CONNS"); min != "" {
		if v, err := strconv.Atoi(min); err == nil {
			cfg.MinConns = int32(v)
		}
	}
	if life := os.Getenv("DB_MAX_CONN_LIFETIME"); life != "" {
		if d, err := time.ParseDuration(life); err == nil {
			cfg.MaxConnLifetime = d
		}
	}
	if idle := os.Getenv("DB_MAX_CONN_IDLE"); idle != "" {
		if d, err := time.ParseDuration(idle); err == nil {
			cfg.MaxConnIdleTime = d
		}
	}

	return Init(ctx, cfg)
}

func Init(ctx context.Context, cfg Config) error {
	once.Do(func() {
		if cfg.DSN == "" {
			initErr = errors.New("db.Init: empty DSN")
			return
		}
		pc, err := pgxpool.ParseConfig(cfg.DSN)
		if err != nil {
			initErr = fmt.Errorf("db.Init: parse config: %w", err)
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
		if initErr != nil {
			return
		}
		ctxPing, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctxPing); err != nil {
			pool.Close()
			pool = nil
			initErr = fmt.Errorf("db.Init: ping: %w", err)
		}
	})
	return initErr
}

func Get() (*pgxpool.Pool, error) {
	if initErr != nil {
		return nil, initErr
	}
	if pool == nil {
		return nil, errors.New("db.Get: pool not initialized")
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
	initErr = nil
	once = sync.Once{}
}

func WithTx(ctx context.Context, opts pgx.TxOptions, fn func(pgx.Tx) error) error {
	p := MustGet()
	tx, err := p.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

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
