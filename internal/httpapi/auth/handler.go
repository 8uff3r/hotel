package auth

import (
	"fmt"
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
	fuego.Get(s, "/profile", au.profileGet)
	fuego.Put(s, "/profile", au.profileUpdate)
	fuego.Put(s, "/profile/password", au.profileChangePassword)
}

type MeResponse struct {
	User        models.SanitizedUser `json:"user"`
	HotelID     string               `json:"hotelId"`
	Permissions []string             `json:"permissions"`
}

func me(c fuego.ContextNoBody) (MeResponse, error) {
	user := c.Value(h.UserKey{}).(models.SanitizedUser)
	hotelID := c.Value(h.HotelIDKey{}).(string)
	permissions := h.GetUserPermissionsFromContext(c)

	return MeResponse{User: user, HotelID: hotelID, Permissions: tokenizePermissions(permissions)}, nil
}

func tokenizePermissions(permissions []models.UserPermissionInfo) []string {
	var permKeys []string
	for _, p := range permissions {
		permKeys = append(permKeys, fmt.Sprintf("%s:%s", p.Page, p.Action))
	}
	if permKeys == nil {
		return []string{}
	}
	return permKeys
}
