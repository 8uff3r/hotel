package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	h "hotel/backend/internal/httpapi"
	system "hotel/backend/internal/httpapi/_system"
	"hotel/backend/internal/httpapi/accounting"
	"hotel/backend/internal/httpapi/auth"
	"hotel/backend/internal/httpapi/guests"
	"hotel/backend/internal/httpapi/parking"
	"hotel/backend/internal/httpapi/reservation"
	"hotel/backend/internal/httpapi/rooms"
	"hotel/backend/internal/httpapi/users"
)

func NewRouter(r *chi.Mux, opts h.Options) http.Handler {

	// Core middlewares, Added inside a group to not interfere with other routes (e.g. SPA routes)
	r.Group(func(r chi.Router) {
		a := &h.API{
			Logger:         opts.Logger,
			Db:             opts.Db,
			SessionCookie:  opts.SessionCookie,
			RequestTimeout: opts.RequestTimeout,
			Services:       opts.Services,
			SessionTTL:     opts.SessionTTL,
		}

		r.Use(a.TimeoutMiddleware)
		r.Use(a.RecoverAndLogMiddleware)

		// Routes

		system.RegisterRoutes(a, r)
		r.Route("/api", func(r chi.Router) {
			auth.RegisterRoutes(a, r)
			r.Group(func(r chi.Router) {
				r.Use(a.Auth)
				accounting.RegisterRoutes(a, r)
				guests.RegisterRoutes(a, r)
				parking.RegisterRoutes(a, r)
				reservation.RegisterRoutes(a, r)
				rooms.RegisterRoutes(a, r)
				users.RegisterRoutes(a, r)
			})
		})
	})

	return r
}
