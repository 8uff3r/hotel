package httpapi

import (
	"context"
	"net/http"
	"time"
)

func (a *API) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), a.requestTimeout)
		defer cancel()
		r = r.WithContext(ctx)
		start := time.Now()
		defer func() {
			if rec := recover(); rec != nil {
				a.logger.Error("panic", "path", r.URL.Path, "err", rec)
				writeErr(w, http.StatusInternalServerError, "internal_error")
			}
			a.logger.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
		}()
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (a *API) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := a.sessionUser(r)
		if err != nil {
			writeErr(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), userKey{}, user)))
	}
}

type userKey struct{}
