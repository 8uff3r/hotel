package httpapi

import (
	"context"
	"errors"
	"hotel/internal/models"
	"net/http"
	"time"
)

// Wraps the request in a timeout context
func (a *API) TimeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.RequestTimeout)
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.sessionUser(r)
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

func (a *API) sessionUser(r *http.Request) (map[string]any, error) {
	cookie, err := r.Cookie(a.SessionCookie)
	if err != nil || cookie.Value == "" {
		return nil, errors.New("missing session")
	}
	user, err := a.Services.Auth.Me(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}
	return SanitizeUser(user), nil
}

func SanitizeUser(u *models.User) map[string]any {
	return map[string]any{"id": u.ID, "email": u.Email, "firstName": u.FirstName, "lastName": u.LastName, "role": u.Role}
}
