package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pintoter/warehouse-api/internal/config"
	"github.com/pintoter/warehouse-api/internal/dbutil/transaction"
	"github.com/pintoter/warehouse-api/internal/migrations"
	productRepository "github.com/pintoter/warehouse-api/internal/repository/product"
	"github.com/pintoter/warehouse-api/internal/server"
	productService "github.com/pintoter/warehouse-api/internal/service/product"
	"github.com/pintoter/warehouse-api/internal/transport"
	"github.com/pintoter/warehouse-api/pkg/database/postgres"
	"github.com/pintoter/warehouse-api/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	ctx := context.Background()

	cfg := config.New()

	syncLogger := initLogger(ctx, cfg)
	defer syncLogger()

	err := migrations.Do(&cfg.DB)
	if err != nil {
		logger.FatalKV(ctx, "Failed init migrations", "err", err)
	}

	db, err := postgres.New(&cfg.DB)
	if err != nil {
		logger.FatalKV(ctx, "Failed connect database", "err", err)
	}
	logger.InfoKV(ctx, "PostgreSQL connected", "Stats", db.Stats())

	repository := productRepository.NewRepository(db)
	txManager := transaction.NewTransactionManager(db)
	service := productService.NewService(repository, txManager)
	handler := transport.NewHandler(service)
	server := server.New(handler, &cfg.HTTP)

	server.Run()
	logger.InfoKV(ctx, "Starting server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	select {
	case <-quit:
		logger.InfoKV(ctx, "Starting gracefully shutdown")
	case err = <-server.Notify():
		logger.FatalKV(ctx, "Failed starting server", "err", err.Error())
	}

	if err := server.Shutdown(); err != nil {
		logger.FatalKV(ctx, "Failed shutdown server", "err", err.Error())
	}
}

type LogConfig interface {
	GetLevel() string
	GetName() string
}

func initLogger(_ context.Context, cfg LogConfig) (syncFn func()) {
	loggingLevel := zap.InfoLevel
	if cfg.GetLevel() == logger.DebugLevel {
		loggingLevel = zap.DebugLevel
	}

	loggerConfig := zap.NewProductionEncoderConfig()

	loggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(loggerConfig),
		os.Stderr,
		zap.NewAtomicLevelAt(loggingLevel),
	)

	notSuggaredLogger := zap.New(consoleCore)

	sugarLogger := notSuggaredLogger.Sugar()
	logger.SetLogger(sugarLogger.With(
		"service", cfg.GetName(),
	))

	return func() {
		_ = notSuggaredLogger.Sync()
	}
}
