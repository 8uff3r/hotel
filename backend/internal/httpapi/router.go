package httpapi

import (
	"log/slog"
	"net/http"
	"time"

	"hotel/backend/internal/models"
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

func NewRouter(opts Options) http.Handler {
	a := &API{logger: opts.Logger, sessionCookie: opts.SessionCookie, requestTimeout: opts.RequestTimeout, services: opts.Services}
	mux := http.NewServeMux()

	a.registerSystemRoutes(mux)
	a.registerAuthRoutes(mux)
	a.registerUserRoutes(mux)
	a.registerRoomRoutes(mux)
	a.registerGuestRoutes(mux)
	a.registerReservationRoutes(mux)
	a.registerAccountingRoutes(mux)
	a.registerParkingRoutes(mux)

	return a.middleware(mux)
}

func (a *API) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", a.health)
	mux.HandleFunc("GET /readyz", a.ready)
}

func (a *API) registerAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", a.login)
	mux.HandleFunc("POST /api/auth/logout", a.auth(a.logout))
	mux.HandleFunc("GET /api/auth/me", a.auth(a.me))
}

func (a *API) registerUserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/users", a.auth(a.usersList))
	mux.HandleFunc("POST /api/users", a.auth(a.usersCreate))
}

func (a *API) registerRoomRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/rooms", a.auth(a.listModel(&models.Room{})))
	mux.HandleFunc("POST /api/rooms", a.auth(a.createModel(&models.Room{})))
	mux.HandleFunc("GET /api/rooms/{id}", a.auth(a.getModel(&models.Room{})))
	mux.HandleFunc("PUT /api/rooms/{id}", a.auth(a.updateModel(&models.Room{})))
	mux.HandleFunc("DELETE /api/rooms/{id}", a.auth(a.deleteModel(&models.Room{})))
}

func (a *API) registerGuestRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/guests", a.auth(a.listModel(&models.Guest{})))
	mux.HandleFunc("POST /api/guests", a.auth(a.createModel(&models.Guest{})))
	mux.HandleFunc("GET /api/guests/{id}", a.auth(a.getModel(&models.Guest{})))
	mux.HandleFunc("PUT /api/guests/{id}", a.auth(a.updateModel(&models.Guest{})))
}

func (a *API) registerReservationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reservations", a.auth(a.listModel(&models.Reservation{})))
	mux.HandleFunc("POST /api/reservations", a.auth(a.createModel(&models.Reservation{})))
	mux.HandleFunc("GET /api/reservations/{id}", a.auth(a.getModel(&models.Reservation{})))
	mux.HandleFunc("PUT /api/reservations/{id}", a.auth(a.updateModel(&models.Reservation{})))
	mux.HandleFunc("POST /api/reservations/{id}/check-in", a.auth(a.reservationsCheckIn))
	mux.HandleFunc("POST /api/reservations/{id}/check-out", a.auth(a.reservationsCheckOut))
}

func (a *API) registerAccountingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/accounts", a.auth(a.listModel(&models.Account{})))
	mux.HandleFunc("POST /api/accounts", a.auth(a.createModel(&models.Account{})))
	mux.HandleFunc("GET /api/expenses", a.auth(a.listModel(&models.Expense{})))
	mux.HandleFunc("POST /api/expenses", a.auth(a.createModel(&models.Expense{})))
	mux.HandleFunc("GET /api/income", a.auth(a.listModel(&models.Income{})))
	mux.HandleFunc("POST /api/income", a.auth(a.createModel(&models.Income{})))
}

func (a *API) registerParkingRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/parking/lots", a.auth(a.listModel(&models.ParkingLot{})))
	mux.HandleFunc("POST /api/parking/lots", a.auth(a.createModel(&models.ParkingLot{})))
	mux.HandleFunc("GET /api/parking/lots/{id}", a.auth(a.getModel(&models.ParkingLot{})))
	mux.HandleFunc("PUT /api/parking/lots/{id}", a.auth(a.updateModel(&models.ParkingLot{})))
	mux.HandleFunc("DELETE /api/parking/lots/{id}", a.auth(a.deleteModel(&models.ParkingLot{})))

	mux.HandleFunc("GET /api/parking/spots", a.auth(a.listModel(&models.ParkingSpot{})))
	mux.HandleFunc("POST /api/parking/spots", a.auth(a.createModel(&models.ParkingSpot{})))
	mux.HandleFunc("GET /api/parking/spots/{id}", a.auth(a.getModel(&models.ParkingSpot{})))
	mux.HandleFunc("PUT /api/parking/spots/{id}", a.auth(a.updateModel(&models.ParkingSpot{})))
	mux.HandleFunc("DELETE /api/parking/spots/{id}", a.auth(a.deleteModel(&models.ParkingSpot{})))

	mux.HandleFunc("GET /api/parking/vehicles", a.auth(a.listModel(&models.Vehicle{})))
	mux.HandleFunc("POST /api/parking/vehicles", a.auth(a.createModel(&models.Vehicle{})))
	mux.HandleFunc("GET /api/parking/vehicles/{id}", a.auth(a.getModel(&models.Vehicle{})))
	mux.HandleFunc("PUT /api/parking/vehicles/{id}", a.auth(a.updateModel(&models.Vehicle{})))
	mux.HandleFunc("DELETE /api/parking/vehicles/{id}", a.auth(a.deleteModel(&models.Vehicle{})))

	mux.HandleFunc("GET /api/parking/transactions", a.auth(a.listModel(&models.ParkingTransaction{})))
	mux.HandleFunc("POST /api/parking/transactions", a.auth(a.createModel(&models.ParkingTransaction{})))
	mux.HandleFunc("GET /api/parking/transactions/{id}", a.auth(a.getModel(&models.ParkingTransaction{})))
	mux.HandleFunc("POST /api/parking/transactions/{id}/check-out", a.auth(a.transactionsCheckOut))
}
