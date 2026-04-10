package auth

import (
	h "hotel/backend/internal/httpapi"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AuthModule struct {
	*h.API
}

func (m AuthModule) RegisterRoutes(api *h.API, r chi.Router) {
	au := AuthModule{API: api}

	r.Post("/login", au.loginHandler)

	r.Group(func(r chi.Router) {
		r.Use(au.Auth)
		r.Post("/logout", au.logout)
		r.Get("/me", me)
	})
}

func me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(h.UserKey{})
	h.WriteJSON(w, 200, map[string]any{"user": user})
}
