package services

import (
	"strings"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type ServicesModule struct {
	*h.API
}

func (m ServicesModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	sm := ServicesModule{api}

	fuego.Get(s, "/", sm.servicesList)
	fuego.Post(s, "/", sm.servicesCreate)
	fuego.Get(s, "/{id}", sm.serviceView)
	fuego.Put(s, "/{id}", sm.serviceUpdate)
	fuego.Delete(s, "/{id}", sm.serviceDelete)
}

func (sm *ServicesModule) servicesList(c fuego.ContextNoBody) (h.PaginatedResponse[models.Service], error) {
	var rows []models.Service
	page := max(c.QueryParamInt("page"), 1)
	limit := c.QueryParamInt("limit")
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	q := sm.Db.WithContext(c).Model(&models.Service{})
	if hotelID := c.QueryParam("hotelId"); hotelID != "" {
		q = q.Where("hotel_id = ?", hotelID)
	}

	var total int64
	q.Count(&total)

	if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return h.PaginatedResponse[models.Service]{}, fuego.InternalServerError{Title: "query_failed"}
	}
	return h.PaginatedResponse[models.Service]{
		Data:       rows,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}

func (sm *ServicesModule) serviceView(c fuego.ContextNoBody) (models.Service, error) {
	var zero models.Service
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	var service models.Service
	if err := sm.Db.First(&service, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	return service, nil
}

func (sm *ServicesModule) servicesCreate(c fuego.ContextWithBody[models.Service]) (models.Service, error) {
	var zero models.Service
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}
	body.Name = strings.TrimSpace(body.Name)
	if err := sm.Db.WithContext(c).Create(&body).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}
	c.SetStatus(201)
	return body, nil
}

func (sm *ServicesModule) serviceUpdate(c fuego.ContextWithBody[models.Service]) (models.Service, error) {
	var zero models.Service
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}
	body.ID = id
	res := sm.Db.WithContext(c).Model(&models.Service{}).Where("id = ?", id).Updates(body)
	if res.Error != nil {
		return zero, fuego.BadRequestError{Title: "update_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	var updated models.Service
	sm.Db.First(&updated, id)
	return updated, nil
}

func (sm *ServicesModule) serviceDelete(c fuego.ContextNoBody) (deleteResponse, error) {
	var zero deleteResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	res := sm.Db.WithContext(c).Delete(&models.Service{}, id)
	if res.Error != nil {
		return zero, fuego.InternalServerError{Title: "delete_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	return deleteResponse{Ok: true}, nil
}

type deleteResponse struct {
	Ok bool `json:"ok"`
}
