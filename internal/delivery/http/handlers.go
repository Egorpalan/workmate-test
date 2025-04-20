package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"

	"github.com/Egorpalan/workmate-test/internal/usecase"
	"github.com/Egorpalan/workmate-test/pkg/logger"
	"go.uber.org/zap"
)

type Handler struct {
	useCase *usecase.UseCase
}

// NewHandler создает новый экземпляр Handler
func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// CreateTask создает новую задачу
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.useCase.Task.CreateTask(r.Context())
	if err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

// GetTask возвращает задачу по ее ID
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "Task ID is required")
		return
	}

	task, err := h.useCase.Task.GetTaskByID(r.Context(), id)
	if err != nil {
		logger.Error("Failed to get task", zap.String("id", id), zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to get task")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

// ListTasks возвращает список задач с пагинацией
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	tasks, err := h.useCase.Task.ListTasks(r.Context(), limit, offset)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		respondWithError(w, http.StatusInternalServerError, "Failed to list tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

// respondWithJSON отправляет JSON-ответ клиенту
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal JSON response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

// respondWithError отправляет сообщение об ошибке клиенту
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
