package auth

import (
	h "hotel/internal/httpapi"

	"github.com/go-fuego/fuego"
)

type AuthModule struct {
	*h.API
}

func RegisterRoutes(api *h.API, s *fuego.Server) {
	au := AuthModule{API: api}

	fuego.Post(s, "/login", au.loginHandler)
	authRequired := fuego.Group(s, "/")
	fuego.Use(authRequired, api.Auth)
	fuego.Post(authRequired, "/logout", au.logout)
	fuego.Get(authRequired, "/me", me)
}

func me(c fuego.ContextNoBody) (map[string]any, error) {
	user := c.Value(h.UserKey{})
	return map[string]any{"user": user}, nil
}
