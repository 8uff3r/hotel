package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	h "hotel/internal/httpapi"
	system "hotel/internal/httpapi/_system"
	"hotel/internal/httpapi/accounting"
	"hotel/internal/httpapi/auth"
	"hotel/internal/httpapi/guests"
	"hotel/internal/httpapi/parking"
	"hotel/internal/httpapi/reservation"
	"hotel/internal/httpapi/rooms"
	"hotel/internal/httpapi/users"
)

func NewRouter(a *h.API, r *chi.Mux) http.Handler {

	// Core middlewares, Added inside a group to not interfere with other routes (e.g. SPA routes)
	r.Group(func(r chi.Router) {
		r.Use(a.TimeoutMiddleware)
		r.Use(a.RecoverAndLogMiddleware)

		// Routes

		SetupRouter(a, r, PathModuleMap{"/": system.SystemModule{}}) // at /healthz and /readyz
		r.Route("/api", func(r chi.Router) {
			SetupRouter(a, r, PathModuleMap{"/auth": auth.AuthModule{}})
			r.Group(func(r chi.Router) {
				r.Use(a.Auth)
				SetupRouter(a, r, PathModuleMap{
					"/accounting":  accounting.AccountingModule{},
					"/guests":      guests.GuestsModule{},
					"/parking":     parking.ParkingModule{},
					"/reservation": reservation.ReservationModule{},
					"/rooms":       rooms.RoomsModule{},
					"/users":       users.UsersModule{},
				})

			})
		})
	})

	return r
}

type ModuleRouter interface {
	RegisterRoutes(api *h.API, r chi.Router)
}

type PathModuleMap map[string]ModuleRouter

func SetupRouter(api *h.API, r chi.Router, modules PathModuleMap) {
	for path, mod := range modules {
		r.Route(path, func(subRouter chi.Router) {
			mod.RegisterRoutes(api, subRouter)
		})
	}
}
