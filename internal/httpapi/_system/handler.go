package system

import (
	h "hotel/internal/httpapi"
	"net/http"

	"github.com/go-fuego/fuego"
)

func RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/healthz", health)
	fuego.Get(s, "/readyz", ready)
}

func health(c fuego.ContextNoBody) (map[string]string, error) {
	return map[string]string{"status": "ok"}
}

func ready(w http.ResponseWriter, _ *http.Request) (map[string]string, error) {
	return map[string]string{"status": "ready"}
}
