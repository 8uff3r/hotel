package accounting

import (
	h "hotel/backend/internal/httpapi"
	"hotel/backend/internal/models"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(api *h.API, r chi.Router) {
	r.Route("/accounting", func(r chi.Router) {
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", api.ListModel(&models.Account{}, nil))
			r.Post("/", api.CreateModel(&models.Account{}))
		})

		r.Route("/expenses", func(r chi.Router) {
			r.Get("/", api.ListModel(&models.Expense{}, nil))
			r.Post("/", api.CreateModel(&models.Expense{}))
		})

		r.Route("/income", func(r chi.Router) {
			r.Get("/", api.ListModel(&models.Income{}, nil))
			r.Post("/", api.CreateModel(&models.Income{}))
		})
	})
}
