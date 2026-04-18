package httpapi

import (
	"context"
	"errors"
	"fmt"
	"hotel/internal/models"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func (a *API) TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.RequestTimeout)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.sessionUser(r)
		fmt.Printf("%s", user)
		if err != nil {
			WriteErr(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), UserKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) RecoverAndLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				a.Logger.Error("panic", "path", r.URL.Path, "err", rec)
				WriteErr(w, http.StatusInternalServerError, "internal_error")
			}

			a.Logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		}()

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type UserKey struct{}

func (a *API) sessionUser(r *http.Request) (models.SanitizedUser, error) {
	var zero models.SanitizedUser
	cookie, err := r.Cookie(a.SessionCookie)
	if err != nil || cookie.Value == "" {
		return zero, errors.New("missing session")
	}

	var s models.Session
	if err := a.Db.WithContext(r.Context()).Preload("User.Roles").Where("id = ? AND expires_at > ?", cookie.Value, time.Now().UTC()).First(&s).Error; err != nil {
		return zero, err
	}
	if !s.User.IsActive {
		return zero, gorm.ErrRecordNotFound
	}
	return SanitizeUser(&s.User), nil
}

func SanitizeUser(u *models.User) models.SanitizedUser {
	return models.SanitizedUser{ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName, Roles: u.Roles}
}
