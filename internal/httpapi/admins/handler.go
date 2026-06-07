package admins

import (
	"strings"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminsModule struct {
	*h.API
}

func (m AdminsModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	am := AdminsModule{api}

	fuego.Get(s, "/", am.adminsList)
	fuego.Post(s, "/", am.adminsCreate)
	fuego.Get(s, "/{id}", am.adminView)
	fuego.Put(s, "/{id}", am.adminUpdate)
	fuego.Delete(s, "/{id}", am.adminDelete)
}

func (am *AdminsModule) adminsList(c fuego.ContextNoBody) (h.PaginatedResponse[models.SanitizedAdmin], error) {
	var rows []models.Admin
	var zero h.PaginatedResponse[models.SanitizedAdmin]

	page := max(c.QueryParamInt("page"), 1)
	limit := c.QueryParamInt("limit")
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	if err := am.Db.WithContext(c).
		Model(&models.Admin{}).
		Order("id DESC").
		Preload("AdminHotels").
		Limit(limit).
		Offset(offset).
		Find(&rows).
		Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	out := make([]models.SanitizedAdmin, 0, len(rows))
	for i := range rows {
		out = append(out, models.SanitizeAdmin(&rows[i]))
	}
	return h.PaginatedResponse[models.SanitizedAdmin]{Data: out}, nil
}

func (am *AdminsModule) adminView(c fuego.ContextNoBody) (models.SanitizedAdmin, error) {
	var zero models.SanitizedAdmin
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var row models.Admin
	if err := am.Db.WithContext(c).Preload("AdminHotels").First(&row, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	return models.SanitizeAdmin(&row), nil
}

type adminCreateDto struct {
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	ContactNumber string   `json:"contactNumber"`
	Email         string   `json:"email"`
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	Role          string   `json:"role"`
	IsSuperAdmin  bool     `json:"isSuperAdmin"`
	HotelIDs      []string `json:"hotelIds"`
}

func (am *AdminsModule) adminsCreate(c fuego.ContextWithBody[adminCreateDto]) (models.SanitizedAdmin, error) {
	var zero models.SanitizedAdmin
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return zero, fuego.InternalServerError{Title: "hash_failed"}
	}

	var adminHotels []models.AdminHotel
	for _, hid := range body.HotelIDs {
		adminHotels = append(adminHotels, models.AdminHotel{HotelID: hid})
	}

	admin := models.Admin{
		FirstName:     strings.TrimSpace(body.FirstName),
		LastName:      strings.TrimSpace(body.LastName),
		ContactNumber: strings.TrimSpace(body.ContactNumber),
		Email:         strings.TrimSpace(body.Email),
		Username:      strings.TrimSpace(body.Username),
		PasswordHash:  string(hash),
		Role:          body.Role,
		IsActive:      true,
		IsSuperAdmin:  body.IsSuperAdmin,
		AdminHotels:   adminHotels,
	}
	if err := am.Db.WithContext(c).Create(&admin).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}
	c.SetStatus(201)
	return models.SanitizeAdmin(&admin), nil
}

type adminUpdateDto struct {
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	ContactNumber string   `json:"contactNumber"`
	Email         string   `json:"email"`
	Username      string   `json:"username"`
	Role          string   `json:"role"`
	IsActive      *bool    `json:"isActive"`
	IsSuperAdmin  *bool    `json:"isSuperAdmin"`
	HotelIDs      []string `json:"hotelIds"`
}

func (am *AdminsModule) adminUpdate(c fuego.ContextWithBody[adminUpdateDto]) (models.SanitizedAdmin, error) {
	var zero models.SanitizedAdmin
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var admin models.Admin
	if err := am.Db.WithContext(c).Preload("AdminHotels").First(&admin, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	if body.FirstName != "" {
		admin.FirstName = strings.TrimSpace(body.FirstName)
	}
	if body.LastName != "" {
		admin.LastName = strings.TrimSpace(body.LastName)
	}
	if body.ContactNumber != "" {
		admin.ContactNumber = strings.TrimSpace(body.ContactNumber)
	}
	if body.Email != "" {
		admin.Email = strings.TrimSpace(body.Email)
	}
	if body.Username != "" {
		admin.Username = strings.TrimSpace(body.Username)
	}
	if body.Role != "" {
		admin.Role = body.Role
	}
	if body.IsActive != nil {
		admin.IsActive = *body.IsActive
	}
	if body.IsSuperAdmin != nil {
		admin.IsSuperAdmin = *body.IsSuperAdmin
	}

	if err := am.Db.WithContext(c).Save(&admin).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "update_failed"}
	}

	// Update hotel associations if provided
	if body.HotelIDs != nil {
		am.Db.Where("admin_id = ?", admin.ID).Delete(&models.AdminHotel{})
		for _, hid := range body.HotelIDs {
			am.Db.Create(&models.AdminHotel{AdminID: admin.ID, HotelID: hid})
		}
	}

	// Reload
	am.Db.Preload("AdminHotels").First(&admin, id)
	return models.SanitizeAdmin(&admin), nil
}

type deleteResponse struct {
	Ok bool `json:"ok"`
}

func (am *AdminsModule) adminDelete(c fuego.ContextNoBody) (deleteResponse, error) {
	var zero deleteResponse
	id, err := h.ParseID(c.PathParam("id"))
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}
	res := am.Db.WithContext(c).Delete(&models.Admin{}, id)
	if res.Error != nil {
		return zero, fuego.InternalServerError{Title: "delete_failed"}
	}
	if res.RowsAffected == 0 {
		return zero, fuego.NotFoundError{}
	}
	return deleteResponse{Ok: true}, nil
}
