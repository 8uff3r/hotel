package auth

import (
	h "hotel/backend/internal/httpapi"
	"net/http"
)

func (a *auth) logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(a.SessionCookie)
	_ = a.Services.Auth.Logout(r.Context(), h.CookieValue(cookie))
	http.SetCookie(w, &http.Cookie{Name: a.SessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}
