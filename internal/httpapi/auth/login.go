package auth

import (
	"context"
	"errors"
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"net/http"
	"strings"
	"time"

	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	User models.SanitizedUser `json:"user"`
}

func (a *AuthModule) loginHandler(c fuego.ContextWithBody[loginDto]) (loginResponse, error) {
	req, err := c.Body()
	var zero loginResponse
	if err != nil {
		return zero, nil
	}
	user, sid, expires, err := a.login(c, req.Email, req.Password)
	if err != nil {
		return zero, fuego.UnauthorizedError{Title: "invalid_credentials"}
	}
	c.SetCookie(http.Cookie{Name: a.SessionCookie, Value: sid, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: expires})
	return loginResponse{User: h.SanitizeUser(user)}, nil
}

func (a *AuthModule) login(ctx context.Context, email, password string) (*models.User, string, time.Time, error) {
	user, err := a.GetUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, "", time.Time{}, errors.New("invalid credentials")
	}
	token, err := h.RandomToken()
	if err != nil {
		return nil, "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(a.SessionTTL)
	if err := a.CreateSession(ctx, models.Session{ID: token, UserID: user.ID, ExpiresAt: expires}); err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expires, nil
}

func (a *AuthModule) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := a.Db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).Preload("Roles").First(&user).Error; err != nil {
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
