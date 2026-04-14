package system

import (
	h "hotel/internal/httpapi"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SystemModule struct{}

func (m SystemModule) RegisterRoutes(api *h.API, r chi.Router) {
	r.Get("/healthz", health)
	r.Get("/readyz", ready)
}

func health(w http.ResponseWriter, _ *http.Request) {
	h.WriteJSON(w, 200, map[string]string{"status": "ok"})
}

func ready(w http.ResponseWriter, _ *http.Request) {
	h.WriteJSON(w, 200, map[string]string{"status": "ready"})
}
