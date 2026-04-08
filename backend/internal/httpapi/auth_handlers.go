package httpapi

import (
	"errors"
	"net/http"

	"hotel/backend/internal/models"
)

func (a *API) login(w http.ResponseWriter, r *http.Request) {
	var req struct{ Email, Password string }
	if !decode(r, &req, w) {
		return
	}
	user, sid, expires, err := a.services.Auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid_credentials")
		return
	}
	http.SetCookie(w, &http.Cookie{Name: a.sessionCookie, Value: sid, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: expires})
	writeJSON(w, 200, map[string]any{"user": sanitizeUser(user)})
}

func (a *API) logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(a.sessionCookie)
	_ = a.services.Auth.Logout(r.Context(), cookieValue(cookie))
	http.SetCookie(w, &http.Cookie{Name: a.sessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	writeJSON(w, 200, map[string]bool{"ok": true})
}

func (a *API) me(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey{})
	writeJSON(w, 200, map[string]any{"user": user})
}

func (a *API) sessionUser(r *http.Request) (map[string]any, error) {
	cookie, err := r.Cookie(a.sessionCookie)
	if err != nil || cookie.Value == "" {
		return nil, errors.New("missing session")
	}
	user, err := a.services.Auth.Me(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}
	return sanitizeUser(user), nil
}

func sanitizeUser(u *models.User) map[string]any {
	return map[string]any{"id": u.ID, "email": u.Email, "firstName": u.FirstName, "lastName": u.LastName, "role": u.Role}
}
