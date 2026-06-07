package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	h "hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	User        models.SanitizedUser `json:"user"`
	HotelID     string               `json:"hotelId"`
	Permissions []string             `json:"permissions"`
}

func (a *AuthModule) loginHandler(c fuego.ContextWithBody[loginDto]) (loginResponse, error) {
	lang := c.Header("Accept-Language")
	if lang == "" {
		lang = "fa"
	}
	req, err := c.Body()
	var zero loginResponse
	if err != nil {
		return zero, nil
	}

	// Try user login first, then admin
	user, sid, expires, err := a.loginUser(c, req.Email, req.Password)
	if err != nil {
		admin, adminSid, adminExpires, adminErr := a.loginAdmin(c, req.Email, req.Password)
		if adminErr != nil {
			return zero, fuego.UnauthorizedError{Title: "invalid_credentials"}
		}
		c.SetCookie(http.Cookie{Name: a.SessionCookie, Value: adminSid, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: adminExpires})

		adminHotels, hotelID := a.getAdminHotels(admin.ID)
		userResponse := models.SanitizedUser{
			ID:            admin.ID,
			FirstName:     admin.FirstName,
			LastName:      admin.LastName,
			Email:         admin.Email,
			Username:      admin.Username,
			ContactNumber: admin.ContactNumber,
			Role:          admin.Role,
			IsAdmin:       true,
			AdminHotels:   adminHotels,
		}
		return loginResponse{User: userResponse, HotelID: hotelID, Permissions: []string{}}, nil
	}

	c.SetCookie(http.Cookie{Name: a.SessionCookie, Value: sid, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: expires})

	userHotels, hotelID := a.getUserHotels(user.ID)

	var roles []models.PermissionTemplate
	for _, v := range user.Roles {
		roles = append(roles, v.Template)
	}
	userResponse := h.SanitizeUser(user, roles)
	userResponse.UserHotels = userHotels

	permissions := a.getUserPermissions(user.ID, lang)

	return loginResponse{User: userResponse, HotelID: hotelID, Permissions: tokenizePermissions(permissions)}, nil
}

func (a *AuthModule) getUserPermissions(userID uint, lang string) []models.UserPermissionInfo {
	var userPerms []models.UserPermission
	if err := a.Db.Preload("Permission").Where("user_id = ?", userID).Find(&userPerms).Error; err != nil {
		return nil
	}

	result := make([]models.UserPermissionInfo, 0, len(userPerms))
	for _, up := range userPerms {
		result = append(result, models.UserPermissionInfo{
			PermissionID: up.PermissionID,
			Page:         up.Permission.Resource,
			Action:       up.Permission.Action,
			Label:        up.Permission.Translation[lang],
			Category:     up.Permission.Category,
			Granted:      up.Granted,
		})
	}
	return result
}

func (a *AuthModule) getUserHotels(userID uint) ([]models.UserHotelInfo, string) {
	var userHotels []models.UserHotel
	if err := a.Db.Preload("Hotel").Where("user_id = ?", userID).Find(&userHotels).Error; err != nil {
		return nil, ""
	}

	result := make([]models.UserHotelInfo, len(userHotels))
	var defaultHotelID string
	for i, uh := range userHotels {
		result[i] = models.UserHotelInfo{
			HotelID: uh.HotelID,
			Hotel:   uh.Hotel,
		}
		if i == 0 {
			defaultHotelID = uh.HotelID
		}
	}
	return result, defaultHotelID
}

func (a *AuthModule) getAdminHotels(adminID uint) ([]models.AdminHotelInfo, string) {
	var adminHotels []models.AdminHotel
	if err := a.Db.Preload("Hotel").Where("admin_id = ?", adminID).Find(&adminHotels).Error; err != nil {
		return nil, ""
	}
	result := make([]models.AdminHotelInfo, len(adminHotels))
	var defaultHotelID string
	for i, ah := range adminHotels {
		result[i] = models.AdminHotelInfo{
			HotelID: ah.HotelID,
			Hotel:   ah.Hotel,
		}
		if i == 0 {
			defaultHotelID = ah.HotelID
		}
	}
	return result, defaultHotelID
}

func (a *AuthModule) loginUser(ctx context.Context, email, password string) (*models.User, string, time.Time, error) {
	user, err := a.GetUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, "", time.Time{}, errors.New("invalid credentials")
	}
	token, err := h.RandomToken()
	if err != nil {
		return nil, "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(a.SessionTTL)
	if err := a.CreateSession(ctx, models.Session{ID: token, UserID: user.ID, ExpiresAt: expires, IsAdmin: false}); err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expires, nil
}

func (a *AuthModule) loginAdmin(ctx context.Context, email, password string) (*models.Admin, string, time.Time, error) {
	var admin models.Admin
	if err := a.Db.WithContext(ctx).Where("email = ? AND is_active = ?", strings.TrimSpace(email), true).First(&admin).Error; err != nil {
		return nil, "", time.Time{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)) != nil {
		return nil, "", time.Time{}, errors.New("invalid credentials")
	}
	token, err := h.RandomToken()
	if err != nil {
		return nil, "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(a.SessionTTL)
	if err := a.CreateSession(ctx, models.Session{ID: token, AdminID: &admin.ID, ExpiresAt: expires, IsAdmin: true}); err != nil {
		return nil, "", time.Time{}, err
	}
	return &admin, token, expires, nil
}

func (a *AuthModule) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := a.Db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).Preload("Roles.Template").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (a *AuthModule) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	var s models.Session
	if err := a.Db.WithContext(ctx).Preload("User").Where("id = ? AND expires_at > ?", sessionID, time.Now().UTC()).First(&s).Error; err != nil {
		return nil, err
	}
	if !s.User.IsActive {
		return nil, gorm.ErrRecordNotFound
	}
	return &s.User, nil
}

func (a *AuthModule) CreateSession(ctx context.Context, session models.Session) error {
	return a.Db.WithContext(ctx).Create(&session).Error
}

func (a *AuthModule) DeleteSession(ctx context.Context, sessionID string) error {
	return a.Db.WithContext(ctx).Delete(&models.Session{}, "id = ?", sessionID).Error
}

func (a *AuthModule) CleanupExpired(ctx context.Context) error {
	return a.Db.WithContext(ctx).Delete(&models.Session{}, "expires_at <= ?", time.Now().UTC()).Error
}
