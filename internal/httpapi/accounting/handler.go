package accounting

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type AccountingModule struct{}

func (m AccountingModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	accountsGroup := fuego.Group(s, "/accounts")

	fuego.Get(accountsGroup, "/", h.ListModel(api.Db, models.Account{}))
	fuego.Post(accountsGroup, "/", h.CreateModel(api.Db, models.Account{}))

	expensesGroup := fuego.Group(s, "/expenses")
	fuego.Get(expensesGroup, "/", h.ListModel(api.Db, models.Expense{}))
	fuego.Post(expensesGroup, "/", h.CreateModel(api.Db, models.Expense{}))

	incomeGroup := fuego.Group(s, "/income")
	fuego.Get(incomeGroup, "/", h.ListModel(api.Db, models.Income{}))
	fuego.Post(incomeGroup, "/", h.CreateModel(api.Db, models.Income{}))
}
