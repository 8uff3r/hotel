package hotels

import (
	"errors"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type HotelsModule struct {
	*h.API
}

func (m HotelsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	fuego.Get(s, "/", h.ListModel[models.Hotel](api.Db))
	fuego.Post(s, "/", hotelsCreate(api.Db))
	fuego.Get(s, "/{id}", h.GetModel[models.Hotel](api.Db))
	fuego.Put(s, "/{id}", h.UpdateModel[models.Hotel](api.Db))
	fuego.Delete(s, "/{id}", hotelsDelete(api.Db))
	fuego.Get(s, "/{id}/settings", hotelSettingsGet(api.Db))
	fuego.Put(s, "/{id}/settings", hotelSettingsUpdate(api.Db))

	fuego.Get(s, "/{id}/pictures", hotelPicturesGet(api.Db))
	fuego.Post(s, "/{id}/pictures", hotelPicturesAdd(api.Db))
	fuego.Delete(s, "/{id}/pictures/{pictureId}", hotelPicturesDelete(api.Db))
}

func hotelsCreate(db *gorm.DB) h.FuegoHandler[models.Hotel, models.Hotel, any] {
	return func(c fuego.ContextWithBody[models.Hotel]) (models.Hotel, error) {
		var zero models.Hotel
		body, err := c.Body()
		if err != nil {
			return zero, fuego.BadRequestError{}
		}
		if body.Code == "" {
			return zero, fuego.BadRequestError{Title: "code_required"}
		}
		if err := db.WithContext(c.Context()).Create(&body).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		c.SetStatus(201)
		return body, nil
	}
}

type deleteResponse struct {
	Ok bool `json:"ok"`
}

func hotelsDelete(db *gorm.DB) func(c fuego.ContextNoBody) (deleteResponse, error) {
	return func(c fuego.ContextNoBody) (deleteResponse, error) {
		var zero deleteResponse
		id := c.PathParam("id")
		if id == "" {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}

		// Check for active users
		var userCount int64
		if err := db.Model(&models.User{}).Where("hotel_id = ?", id).Count(&userCount).Error; err != nil {
			return zero, fuego.InternalServerError{Title: "check_failed"}
		}
		if userCount > 0 {
			return zero, fuego.BadRequestError{Title: "hotel_has_active_users"}
		}

		// Check for rooms
		var roomCount int64
		if err := db.Model(&models.Room{}).Where("hotel_id = ?", id).Count(&roomCount).Error; err != nil {
			return zero, fuego.InternalServerError{Title: "check_failed"}
		}
		if roomCount > 0 {
			return zero, fuego.BadRequestError{Title: "hotel_has_rooms"}
		}

		// Check for stays
		var stayCount int64
		if err := db.Model(&models.Stay{}).Where("hotel_id = ?", id).Count(&stayCount).Error; err != nil {
			return zero, fuego.InternalServerError{Title: "check_failed"}
		}
		if stayCount > 0 {
			return zero, fuego.BadRequestError{Title: "hotel_has_stays"}
		}

		res := db.WithContext(c).Delete(&models.Hotel{}, "id = ?", id)
		if res.Error != nil {
			return zero, fuego.InternalServerError{Title: "delete_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return deleteResponse{Ok: true}, nil
	}
}

type HotelSettingResponse struct {
	Setting *models.HotelSetting `json:"setting"`
}

func hotelSettingsGet(db *gorm.DB) h.FuegoAnyHandler[HotelSettingResponse] {
	return func(c fuego.ContextNoBody) (HotelSettingResponse, error) {
		id := c.PathParam("id")
		if id == "" {
			return HotelSettingResponse{}, fuego.BadRequestError{Title: "invalid_id"}
		}
		var setting models.HotelSetting
		if err := db.Where("hotel_id = ?", id).First(&setting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return HotelSettingResponse{Setting: nil}, nil
			}
			return HotelSettingResponse{}, fuego.InternalServerError{Title: "query_failed"}
		}
		return HotelSettingResponse{Setting: &setting}, nil
	}
}

func hotelSettingsUpdate(db *gorm.DB) h.FuegoHandler[models.HotelSetting, models.HotelSetting, any] {
	return func(c fuego.ContextWithBody[models.HotelSetting]) (models.HotelSetting, error) {
		var zero models.HotelSetting
		id := c.PathParam("id")
		if id == "" {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		body, err := c.Body()
		if err != nil {
			return zero, fuego.BadRequestError{}
		}
		body.HotelID = id

		var existing models.HotelSetting
		if err := db.Where("hotel_id = ?", id).First(&existing).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err := db.Create(&body).Error; err != nil {
					return zero, fuego.BadRequestError{Title: "create_failed"}
				}
				return body, nil
			}
			return zero, fuego.InternalServerError{Title: "query_failed"}
		}

		body.ID = existing.ID
		if err := db.Save(&body).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "update_failed"}
		}
		return body, nil
	}
}

type hotelPicturesResponse struct {
	Data []models.HotelPicture `json:"data"`
}

func hotelPicturesGet(db *gorm.DB) func(c fuego.ContextNoBody) (hotelPicturesResponse, error) {
	return func(c fuego.ContextNoBody) (hotelPicturesResponse, error) {
		id := c.PathParam("id")
		if id == "" {
			return hotelPicturesResponse{}, fuego.BadRequestError{Title: "invalid_id"}
		}
		var pictures []models.HotelPicture
		if err := db.Where("hotel_id = ?", id).Find(&pictures).Error; err != nil {
			return hotelPicturesResponse{}, fuego.InternalServerError{Title: "query_failed"}
		}
		return hotelPicturesResponse{Data: pictures}, nil
	}
}

type hotelPictureDto struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

func hotelPicturesAdd(db *gorm.DB) func(c fuego.ContextWithBody[hotelPictureDto]) (models.HotelPicture, error) {
	return func(c fuego.ContextWithBody[hotelPictureDto]) (models.HotelPicture, error) {
		var zero models.HotelPicture
		id := c.PathParam("id")
		if id == "" {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		body, err := c.Body()
		if err != nil {
			return zero, fuego.BadRequestError{}
		}
		if body.URL == "" {
			return zero, fuego.BadRequestError{Title: "url_required"}
		}
		picture := models.HotelPicture{
			HotelID:     id,
			URL:         body.URL,
			Description: body.Description,
		}
		if err := db.Create(&picture).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		return picture, nil
	}
}

func hotelPicturesDelete(db *gorm.DB) func(c fuego.ContextNoBody) (deleteResponse, error) {
	return func(c fuego.ContextNoBody) (deleteResponse, error) {
		var zero deleteResponse
		pictureID := c.PathParam("pictureId")
		if pictureID == "" {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		res := db.Delete(&models.HotelPicture{}, pictureID)
		if res.Error != nil {
			return zero, fuego.InternalServerError{Title: "delete_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return deleteResponse{Ok: true}, nil
	}
}
