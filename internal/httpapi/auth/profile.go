package auth

import (
	"strings"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type profileUpdateDto struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type changePasswordDto struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (a *AuthModule) profileGet(c fuego.ContextNoBody) (models.SanitizedUser, error) {
	user := c.Value(h.UserKey{}).(models.SanitizedUser)
	return user, nil
}

func (a *AuthModule) profileUpdate(c fuego.ContextWithBody[profileUpdateDto]) (models.SanitizedUser, error) {
	var zero models.SanitizedUser
	sanitized := c.Value(h.UserKey{}).(models.SanitizedUser)

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	var row models.User
	if err := a.Db.WithContext(c).
		Model(&models.User{}).
		Preload("UserHotels").
		Preload("Roles.Template").
		First(&row, sanitized.ID).
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

	if err := a.Db.WithContext(c).Save(&row).Error; err != nil {
		return zero, fuego.BadRequestError{Title: "update_failed"}
	}

	var roles []models.PermissionTemplate
	for _, v := range row.Roles {
		roles = append(roles, v.Template)
	}

	return h.SanitizeUser(&row, roles), nil
}

func (a *AuthModule) profileChangePassword(c fuego.ContextWithBody[changePasswordDto]) (map[string]string, error) {
	var zero map[string]string
	sanitized := c.Value(h.UserKey{}).(models.SanitizedUser)

	body, err := c.Body()
	if err != nil {
		return zero, fuego.BadRequestError{}
	}

	if body.CurrentPassword == "" || body.NewPassword == "" {
		return zero, fuego.BadRequestError{Title: "missing_fields"}
	}

	if len(body.NewPassword) < 6 {
		return zero, fuego.BadRequestError{Title: "password_too_short"}
	}

	var row models.User
	if err := a.Db.WithContext(c).Model(&models.User{}).First(&row, sanitized.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return zero, fuego.NotFoundError{Title: "not_found"}
		}
		return zero, fuego.InternalServerError{Title: "query_failed"}
	}

	if bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(body.CurrentPassword)) != nil {
		return zero, fuego.UnauthorizedError{Title: "invalid_current_password"}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return zero, fuego.InternalServerError{Title: "hash_failed"}
	}

	if err := a.Db.WithContext(c).Model(&row).Update("password_hash", string(hash)).Error; err != nil {
		return zero, fuego.InternalServerError{Title: "update_failed"}
	}

	return map[string]string{"message": "password_updated"}, nil
}
