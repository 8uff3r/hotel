package auth

import (
	h "hotel/internal/httpapi"
	"hotel/internal/models"
	"net/http"

	"github.com/go-fuego/fuego"
)

func (a *AuthModule) logout(c fuego.ContextNoBody) (map[string]bool, error) {
	cookie, _ := c.Cookie(a.SessionCookie)
	_ = a.Db.WithContext(c).Delete(&models.Session{}, "id = ?", h.CookieValue(cookie)).Error
	c.SetCookie(http.Cookie{Name: a.SessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true})
	return map[string]bool{"ok": true}, nil
}
