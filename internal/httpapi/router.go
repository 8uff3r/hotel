package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"hotel/backend/internal/models"
	"hotel/backend/internal/repository"
	"hotel/backend/internal/service"
)

type Options struct {
	Logger         *slog.Logger
	SessionCookie  string
	RequestTimeout time.Duration
	Services       service.Services
}

type API struct {
	logger         *slog.Logger
	sessionCookie  string
	requestTimeout time.Duration
	services       service.Services
}

func NewRouter(r *chi.Mux, opts Options) http.Handler {
	a := &API{
		logger:         opts.Logger,
		sessionCookie:  opts.SessionCookie,
		requestTimeout: opts.RequestTimeout,
		services:       opts.Services,
	}

	// Core middlewares
	r.Group(func(r chi.Router) {
		r.Use(a.timeoutMiddleware)
		r.Use(a.recoverAndLogMiddleware)

		// Routes

		a.registerSystemRoutes(r)
		r.Route("/api", func(r chi.Router) {
			a.registerAuthRoutes(r)
			r.Group(func(r chi.Router) {
				r.Use(a.auth)
				a.registerUserRoutes(r)
				a.registerRoomRoutes(r)
				a.registerGuestRoutes(r)
				a.registerReservationRoutes(r)
				a.registerAccountingRoutes(r)
				a.registerParkingRoutes(r)
			})
		})
	})

	return r
}

func (a *API) registerSystemRoutes(r chi.Router) {
	r.Get("/healthz", a.health)
	r.Get("/readyz", a.ready)
}

func (a *API) registerAuthRoutes(r chi.Router) {
	r.Post("/auth/login", a.login)

	r.With(a.auth).Route("/auth", func(r chi.Router) {
		r.Post("/logout", a.logout)
		r.Get("/me", a.me)
	})
}

func (a *API) registerUserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Get("/", a.usersList)
		r.Post("/", a.usersCreate)
	})
}

func (a *API) registerRoomRoutes(r chi.Router) {
	r.Route("/rooms", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Room{}, &repository.ListOptions{
			Preload: []string{"Amenities"},
		}))
		r.Post("/", a.createModel(&models.Room{}))
		r.Get("/{id}", a.getModel(&models.Room{}))
		r.Put("/{id}", a.updateModel(&models.Room{}))
		r.Delete("/{id}", a.deleteModel(&models.Room{}))

		r.Get("/amenities", a.listModel(&models.Amenity{}, nil))
	})
}

func (a *API) registerGuestRoutes(r chi.Router) {
	r.Route("/guests", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Guest{}, nil))
		r.Post("/", a.createModel(&models.Guest{}))
		r.Get("/{id}", a.getModel(&models.Guest{}))
		r.Put("/{id}", a.updateModel(&models.Guest{}))
	})
}

func (a *API) registerReservationRoutes(r chi.Router) {
	r.Route("/reservations", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Reservation{}, nil))
		r.Post("/", a.createModel(&models.Reservation{}))
		r.Get("/{id}", a.getModel(&models.Reservation{}))
		r.Put("/{id}", a.updateModel(&models.Reservation{}))

		r.Post("/{id}/check-in", a.reservationsCheckIn)
		r.Post("/{id}/check-out", a.reservationsCheckOut)
	})
}

func (a *API) registerAccountingRoutes(r chi.Router) {
	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Account{}, nil))
		r.Post("/", a.createModel(&models.Account{}))
	})

	r.Route("/expenses", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Expense{}, nil))
		r.Post("/", a.createModel(&models.Expense{}))
	})

	r.Route("/income", func(r chi.Router) {
		r.Get("/", a.listModel(&models.Income{}, nil))
		r.Post("/", a.createModel(&models.Income{}))
	})
}

func (a *API) registerParkingRoutes(r chi.Router) {
	r.Route("/parking", func(r chi.Router) {

		r.Route("/lots", func(r chi.Router) {
			r.Get("/", a.listModel(&models.ParkingLot{}, nil))
			r.Post("/", a.createModel(&models.ParkingLot{}))
			r.Get("/{id}", a.getModel(&models.ParkingLot{}))
			r.Put("/{id}", a.updateModel(&models.ParkingLot{}))
			r.Delete("/{id}", a.deleteModel(&models.ParkingLot{}))
		})

		r.Route("/spots", func(r chi.Router) {
			r.Get("/", a.listModel(&models.ParkingSpot{}, nil))
			r.Post("/", a.createModel(&models.ParkingSpot{}))
			r.Get("/{id}", a.getModel(&models.ParkingSpot{}))
			r.Put("/{id}", a.updateModel(&models.ParkingSpot{}))
			r.Delete("/{id}", a.deleteModel(&models.ParkingSpot{}))
		})

		r.Route("/vehicles", func(r chi.Router) {
			r.Get("/", a.listModel(&models.Vehicle{}, nil))
			r.Post("/", a.createModel(&models.Vehicle{}))
			r.Get("/{id}", a.getModel(&models.Vehicle{}))
			r.Put("/{id}", a.updateModel(&models.Vehicle{}))
			r.Delete("/{id}", a.deleteModel(&models.Vehicle{}))
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", a.listModel(&models.ParkingTransaction{}, nil))
			r.Post("/", a.createModel(&models.ParkingTransaction{}))
			r.Get("/{id}", a.getModel(&models.ParkingTransaction{}))
			r.Post("/{id}/check-out", a.transactionsCheckOut)
		})
	})
}
