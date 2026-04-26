package auth

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type AuthModule struct {
	*h.API
}

func (m AuthModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	au := AuthModule{API: api}

	fuego.Post(s, "/login", au.loginHandler)
	fuego.Use(s, api.AuthMiddleware)
	fuego.Post(s, "/logout", au.logout)
	fuego.Get(s, "/me", me)
}

type MeResponse struct {
	User    models.SanitizedUser `json:"user"`
	HotelID string               `json:"hotelId"`
}

func me(c fuego.ContextNoBody) (MeResponse, error) {
	user := c.Value(h.UserKey{}).(models.SanitizedUser)
	hotelID := c.Value(h.HotelIDKey{}).(string)

	return MeResponse{User: user, HotelID: hotelID}, nil
}
