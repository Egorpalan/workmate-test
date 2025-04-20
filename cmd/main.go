package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	net "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Egorpalan/workmate-test/config"
	"github.com/Egorpalan/workmate-test/internal/delivery/http"
	"github.com/Egorpalan/workmate-test/internal/repository/postgresql"
	"github.com/Egorpalan/workmate-test/internal/usecase"
	"github.com/Egorpalan/workmate-test/pkg/db"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger.Setup()
	defer logger.GetLogger().Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	dbConn, err := db.NewPostgresConnection(&cfg.DB)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func(dbConn *sqlx.DB) {
		err := dbConn.Close()
		if err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}(dbConn)

	if err := dbConn.PingContext(context.Background()); err != nil {
		logger.Fatal("Database ping failed", zap.Error(err))
	}

	taskRepo := postgresql.NewTaskRepository(dbConn)

	processTask := func(ctx context.Context) (json.RawMessage, error) {
		logger.Info("Starting long running task")
		time.Sleep(3 * time.Minute)

		result := map[string]interface{}{
			"message":   "Task completed successfully",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal result: %w", err)
		}

		logger.Info("Long running task completed")
		return resultJSON, nil
	}

	taskUseCase := usecase.NewTaskUseCase(taskRepo, processTask)
	uc := usecase.NewUseCase(taskUseCase)

	server := http.NewServer(cfg, uc)

	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, net.ErrServerClosed) {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("port", cfg.Server.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown error", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
