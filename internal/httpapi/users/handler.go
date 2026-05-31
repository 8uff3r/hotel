package users

import (
	"strconv"
	"strings"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersModule struct {
	*h.API
}

func (m UsersModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	u := UsersModule{api}

	fuego.Get(s, "/", u.usersList)
	fuego.Post(s, "/", u.usersCreate)
	fuego.Get(s, "/{id}", u.userView)
	fuego.Put(s, "/{id}", u.userUpdate)
}

type userListResponse struct {
	Data []models.SanitizedUser `json:"data"`
}

func (u UsersModule) userView(c fuego.ContextNoBody) (models.SanitizedUser, error) {
	var zero models.SanitizedUser
	id := c.PathParam("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	var row models.User
	if err := u.Db.WithContext(c).
		Model(&models.User{}).
		Preload("UserHotels").
		Preload("Roles.Template").
		First(&row, uint(uid)).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	var roles []models.PermissionTemplate
	for _, v := range row.Roles {
		roles = append(roles, v.Template)
	}
	return h.SanitizeUser(&row, roles), nil
}

func (u *UsersModule) userUpdate(c fuego.ContextWithBody[userUpdateDto]) (models.SanitizedUser, error) {
	var zero models.SanitizedUser
	id := c.PathParam("id")
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return zero, fuego.BadRequestError{Title: "invalid_id"}
	}

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var row models.User
	if err := u.Db.WithContext(c).
		Model(&models.User{}).
		Preload("UserHotels").
		Preload("Roles.Template").
		First(&row, uint(uid)).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	if body.Email != "" {
		row.Email = strings.TrimSpace(body.Email)
	}
	if body.FirstName != "" {
		row.FirstName = strings.TrimSpace(body.FirstName)
	}
	if body.LastName != "" {
		row.LastName = strings.TrimSpace(body.LastName)
	}

	if err := u.Db.WithContext(c).Save(&row).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "update_failed"}
	}
	var roles []models.PermissionTemplate
	for _, v := range row.Roles {
		roles = append(roles, v.Template)
	}

	return h.SanitizeUser(&row, roles), nil
}

func (u *UsersModule) usersList(c fuego.ContextNoBody) (h.PaginatedResponse[models.SanitizedUser], error) {
	var rows []models.User
	var zero h.PaginatedResponse[models.SanitizedUser]

	page := max(c.QueryParamInt("page"), 1)
	limit := c.QueryParamInt("limit")
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	if err := u.Db.WithContext(c).
		Model(&models.User{}).
		Order("id DESC").
		Preload("UserHotels").
		Preload("Roles.Template").
		Limit(limit).
		Offset(offset).
		Find(&rows).
		Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	out := make([]models.SanitizedUser, 0, len(rows))
	for i := range rows {

		var roles []models.PermissionTemplate
		for _, v := range rows[i].Roles {
			roles = append(roles, v.Template)
		}
		out = append(out, h.SanitizeUser(&rows[i], roles))
	}
	return h.PaginatedResponse[models.SanitizedUser]{Data: out}, nil
}

type userCreateDto struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	RoleIDs   *[]uint `json:"roleIds"`
}
type userCreateResponse struct {
	ID uint `json:"id"`
}

type userUpdateDto struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (u *UsersModule) usersCreate(c fuego.ContextWithBody[userCreateDto]) (userCreateResponse, error) {
	var zero userCreateResponse
	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return zero, fuego.InternalServerError{Title: "hash_failed"}
	}

	var roles []models.UserTemplate
	for _, v := range *body.RoleIDs {
		roles = append(roles, models.UserTemplate{TemplateID: v})
	}
	user := &models.User{
		Email:        strings.TrimSpace(body.Email),
		PasswordHash: string(hash),
		FirstName:    strings.TrimSpace(body.FirstName),
		LastName:     strings.TrimSpace(body.LastName),
		Roles:        roles,
		IsActive:     true,
	}
	if err := u.Db.WithContext(c).Create(user).Error; err != nil {
		return userCreateResponse{}, err
	}
	c.SetStatus(201)
	return userCreateResponse{ID: user.ID}, nil
}
