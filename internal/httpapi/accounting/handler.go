package accounting

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-chi/chi/v5"
)

type AccountingModule struct{}

func (m AccountingModule) RegisterRoutes(api *h.API, r chi.Router) {
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
}
