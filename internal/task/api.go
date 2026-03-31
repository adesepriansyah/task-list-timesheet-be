package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/auth"
	"github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
	"github.com/go-chi/chi/v5"
)

type resource struct {
	service Service
}

// RegisterHandlers registers the task handlers.
func RegisterHandlers(r chi.Router, service Service, jwtSecret []byte) {
	res := &resource{service}

	r.Route("/api/tasks", func(r chi.Router) {
		r.Use(auth.JWTMiddleware(jwtSecret))
		r.Post("/", res.create)
		r.Get("/", res.list)
		r.Put("/{id}", res.update)
		r.Delete("/{id}", res.delete)
	})
}

func (res *resource) create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	if err := res.service.CreateTask(r.Context(), req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"data": "Ok"})
}

func (res *resource) list(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	filter := repository.TaskFilter{
		UserID:   userID,
		Search:   r.URL.Query().Get("search"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	tasks, err := res.service.GetTasks(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": tasks})
}

func (res *resource) update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	userID, _ := auth.UserIDFromContext(r.Context())

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	if err := res.service.UpdateTask(r.Context(), id, userID, req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "Ok"})
}

func (res *resource) delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}

	if err := res.service.DeleteTask(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "Ok"})
}
