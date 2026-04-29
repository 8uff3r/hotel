package users

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"strings"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
)

type UsersModule struct {
	*h.API
}

func (m UsersModule) RegisterRoutes(api *h.API, s *fuego.Server) {
	u := UsersModule{api}

	fuego.Get(s, "/", u.usersList)
	fuego.Post(s, "/", u.usersCreate)
}

type userListResponse struct {
	Data []models.SanitizedUser `json:"data"`
}

func (u *UsersModule) usersList(c fuego.ContextNoBody) (h.PaginatedResponse[models.SanitizedUser], error) {
	var rows []models.User
	var zero h.PaginatedResponse[models.SanitizedUser]

	page := c.QueryParamInt("page")
	if page < 1 {
		page = 1
	}
	limit := c.QueryParamInt("limit")
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	if err := u.Db.WithContext(c).Model(&models.User{}).Order("id DESC").Preload("UserHotels").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}
	out := make([]models.SanitizedUser, 0, len(rows))
	for i := range rows {
		out = append(out, h.SanitizeUser(&rows[i]))
	}
	return h.PaginatedResponse[models.SanitizedUser]{Data: out}, nil
}

type userCreateDto struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
type userCreateResponse struct {
	id uint
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
	user := &models.User{
		Email:        strings.TrimSpace(body.Email),
		PasswordHash: string(hash),
		FirstName:    strings.TrimSpace(body.FirstName),
		LastName:     strings.TrimSpace(body.LastName),
		IsActive:     true,
	}
	if err := u.Db.WithContext(c).Create(user).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "create_failed"}
	}
	c.SetStatus(201)
	return userCreateResponse{id: user.ID}, nil
}
