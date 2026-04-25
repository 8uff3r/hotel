package hotels

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
)

type HotelsModule struct {
	*h.API
}

type UserHotelResponse struct {
	UserHotels []models.UserHotelInfo `json:"userHotels"`
	HotelID    string                 `json:"hotelId"`
}

func (m HotelsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	hm := HotelsModule{api}

	fuego.Get(s, "/", h.ListModel(api.Db, models.Hotel{}))
	fuego.Post(s, "/", h.CreateModel(api.Db, models.Hotel{}))
	fuego.Get(s, "/{id}", h.GetModel(api.Db, models.Hotel{}))
	fuego.Put(s, "/{id}", h.UpdateModel(api.Db, models.Hotel{}))
	fuego.Delete(s, "/{id}", h.DeleteModel(api.Db, models.Hotel{}))

	fuego.Get(s, "/my", hm.getUserHotels)
}

func (hm *HotelsModule) getUserHotels(c fuego.ContextNoBody) (UserHotelResponse, error) {
	userVal := c.Context().Value(h.UserKey{})
	if userVal == nil {
		return UserHotelResponse{}, fuego.NotFoundError{}
	}
	user := userVal.(models.SanitizedUser)

	var userHotels []models.UserHotel
	if err := hm.Db.Preload("Hotel").Preload("Role").Where("user_id = ?", user.ID).Find(&userHotels).Error; err != nil {
		return UserHotelResponse{}, err
	}

	result := make([]models.UserHotelInfo, len(userHotels))
	for i, uh := range userHotels {
		result[i] = models.UserHotelInfo{
			HotelID: uh.HotelID,
			Hotel:   uh.Hotel,
			RoleID:  uh.RoleID,
			Role:    uh.Role,
		}
	}

	selectedHotelID := hm.MustGetHotelIDFromCookie(c.Request())

	return UserHotelResponse{
		UserHotels: result,
		HotelID:    selectedHotelID,
	}, nil
}
