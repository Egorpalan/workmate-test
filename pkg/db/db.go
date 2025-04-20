package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Egorpalan/workmate-test/config"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// ConnectWithRetry устанавливает соединение с БД с повторными попытками
func ConnectWithRetry(cfg *config.DBConfig, maxRetries int, retryInterval time.Duration) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	dsn := cfg.GetDSN()

	for i := 0; i < maxRetries; i++ {
		logger.Info("Attempting to connect to database",
			zap.String("attempt", fmt.Sprintf("%d/%d", i+1, maxRetries)),
			zap.String("dsn", dsn))

		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			logger.Info("Successfully connected to database")
			return db, nil
		}

		logger.Error("Failed to connect to database",
			zap.Error(err),
			zap.Duration("retry_interval", retryInterval))

		time.Sleep(retryInterval)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// NewPostgresConnection создает новое подключение к PostgreSQL
func NewPostgresConnection(cfg *config.DBConfig) (*sqlx.DB, error) {
	const (
		defaultMaxRetries    = 5
		defaultRetryInterval = 2 * time.Second
	)

	return ConnectWithRetry(cfg, defaultMaxRetries, defaultRetryInterval)
}

// PingDatabase проверяет доступность БД
func PingDatabase(ctx context.Context, db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
