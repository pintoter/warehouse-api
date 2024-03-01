package postgres

import (
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	driverName = "pgx"
)

type Config interface {
	GetDSN() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxIdleTime() time.Duration
	GetConnMaxLifetime() time.Duration
}

func New(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(driverName, cfg.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "opening db")
	}

	db.SetMaxOpenConns(cfg.GetMaxOpenConns())
	db.SetMaxIdleConns(cfg.GetMaxIdleConns())
	db.SetConnMaxIdleTime(cfg.GetConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "ping DB")
	}

	return db, nil
}
