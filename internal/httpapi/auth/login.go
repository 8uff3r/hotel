package auth

import (
	"context"
	"errors"
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (a *AuthModule) loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if err := h.Decode(r, &req, w); err != nil {
		return
	}
	user, sid, expires, err := a.login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.WriteErr(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: a.SessionCookie, Value: sid, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: expires})
	h.WriteJSON(w, 200, map[string]any{"user": h.SanitizeUser(user)})
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
	if err := a.Db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
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
