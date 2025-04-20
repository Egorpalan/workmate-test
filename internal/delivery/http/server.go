package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Egorpalan/workmate-test/config"
	"github.com/Egorpalan/workmate-test/internal/usecase"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	handler    *Handler
}

// NewServer создает новый экземпляр Server
func NewServer(cfg *config.Config, useCase *usecase.UseCase) *Server {
	handler := NewHandler(useCase)

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
			Handler:      setupRouter(handler),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		handler: handler,
	}
}

// setupRouter настраивает маршруты сервера
func setupRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Route("/api", func(r chi.Router) {
		r.Route("/tasks", func(r chi.Router) {
			r.Post("/", h.CreateTask)
			r.Get("/{id}", h.GetTask)
			r.Get("/", h.ListTasks)
		})
	})

	return r
}

// Run запускает HTTP-сервер
func (s *Server) Run() error {
	logger.Info("Starting HTTP server", zap.String("address", s.httpServer.Addr))
	return s.httpServer.ListenAndServe()
}

// Shutdown останавливает HTTP-сервер
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
