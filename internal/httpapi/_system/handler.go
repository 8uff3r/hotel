package system

import (
	h "hotel/internal/httpapi"

	"github.com/go-fuego/fuego"
)

type SystemModule struct{}

func (m SystemModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/healthz", health)
	fuego.Get(s, "/readyz", ready)
}

func health(c fuego.ContextNoBody) (map[string]string, error) {
	return map[string]string{"status": "ok"}, nil
}

func ready(c fuego.ContextNoBody) (map[string]string, error) {
	return map[string]string{"status": "ready"}, nil
}
