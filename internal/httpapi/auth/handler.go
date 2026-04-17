package auth

import (
	h "hotel/internal/httpapi"

	"github.com/go-fuego/fuego"
)

type AuthModule struct {
	*h.API
}

func (m AuthModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	au := AuthModule{API: api}

	fuego.Post(s, "/login", au.loginHandler)
	fuego.Use(s, api.Auth)
	fuego.Post(s, "/logout", au.logout)
	fuego.Get(s, "/me", me)
}

type response struct {
	User any `json:"user"`
}

func me(c fuego.ContextNoBody) (response, error) {
	user := c.Value(h.UserKey{})
	return response{User: user}, nil
}
