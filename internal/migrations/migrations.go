package migrations

import (
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pintoter/warehouse-api/pkg/logger"
)

const (
	sourceURL = "file://migrations"
)

type Config interface {
	GetDSN() string
}

func Do(cfg Config) error {
	logger.DebugKV(context.Background(), "res", "cfg.DSN()", cfg.GetDSN())
	m, err := migrate.New(sourceURL, cfg.GetDSN())
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.DebugKV(context.Background(), "res", "err", err)
		return err
	}
	defer func() {
		m.Close()
	}()

	return nil
}
