package accounting

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type AccountingModule struct{}

func (m AccountingModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	accountsGroup := fuego.Group(s, "/accounts")

	fuego.Get(accountsGroup, "/", h.ListModel[models.Account](api.Db))
	fuego.Post(accountsGroup, "/", h.CreateModel[models.Account](api.Db))

	expensesGroup := fuego.Group(s, "/expenses")
	fuego.Get(expensesGroup, "/", h.ListModel[models.Expense](api.Db))
	fuego.Post(expensesGroup, "/", h.CreateModel[models.Expense](api.Db))

	incomeGroup := fuego.Group(s, "/income")
	fuego.Get(incomeGroup, "/", h.ListModel[models.Income](api.Db))
	fuego.Post(incomeGroup, "/", h.CreateModel[models.Income](api.Db))

	paymentMethodsGroup := fuego.Group(s, "/payment-methods")
	fuego.Get(paymentMethodsGroup, "/", h.ListModel[models.PaymentMethod](api.Db, h.WithTranslation[models.PaymentMethod]()))

	paymentStatusesGroup := fuego.Group(s, "/payment-statuses")
	fuego.Get(paymentStatusesGroup, "/", h.ListModel[models.PaymentStatus](api.Db, h.WithTranslation[models.PaymentStatus]()))
}
