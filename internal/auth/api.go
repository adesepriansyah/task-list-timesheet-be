package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type resource struct {
	service Service
}

// RegisterHandlers registers the auth handlers.
func RegisterHandlers(r chi.Router, service Service) {
	res := &resource{service}

	r.Route("/api/users", func(r chi.Router) {
		r.Post("/login", res.login)
		r.Post("/logout", res.logout)
	})
}

func (res *resource) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := res.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "unauthorized" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]string{
			"token": token,
		},
	})
}

func (res *resource) logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	// Expected: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	token := parts[1]
	if err := res.service.Logout(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"data": "success",
	})
}
