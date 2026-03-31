package healthcheck

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// RegisterHandlers registers the healthcheck handlers.
func RegisterHandlers(r chi.Router) {
	r.Get("/health", health)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
