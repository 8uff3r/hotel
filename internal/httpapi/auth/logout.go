package auth

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"net/http"
)

func (a *AuthModule) logout(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie(a.SessionCookie)
	_ = a.Db.WithContext(r.Context()).Delete(&models.Session{}, "id = ?", h.CookieValue(cookie)).Error
	http.SetCookie(w, &http.Cookie{Name: a.SessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	h.WriteJSON(w, 200, map[string]bool{"ok": true})
}
