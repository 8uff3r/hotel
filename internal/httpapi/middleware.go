package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"hotel/internal/models"

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
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = "fa"
		}

		userHotels := a.getUserHotelsFromDB(user.ID)
		hotelID := a.resolveHotelID(r, userHotels)
		permissions := a.getUserPermissionsFromDB(user.ID, lang)

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey{}, user)
		ctx = context.WithValue(ctx, UserHotelsKey{}, userHotels)
		ctx = context.WithValue(ctx, HotelIDKey{}, hotelID)
		ctx = context.WithValue(ctx, UserPermissionsKey{}, permissions)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) getUserHotelsFromDB(userID uint) []models.UserHotelInfo {
	var userHotels []models.UserHotel
	if err := a.Db.Preload("Hotel").Where("user_id = ?", userID).Find(&userHotels).Error; err != nil {
		return nil
	}

	result := make([]models.UserHotelInfo, 0, len(userHotels))
	for _, uh := range userHotels {
		result = append(result, models.UserHotelInfo{
			HotelID: uh.HotelID,
			Hotel:   uh.Hotel,
		})
	}
	return result
}

func (a *API) resolveHotelID(r *http.Request, userHotels []models.UserHotelInfo) string {
	cookie, err := r.Cookie(a.HotelCookie)
	if err == nil && cookie.Value != "" {
		for _, uh := range userHotels {
			if uh.HotelID == cookie.Value {
				return cookie.Value
			}
		}
	}

	if len(userHotels) > 0 {
		return userHotels[0].HotelID
	}
	return ""
}

func (a *API) RecoverAndLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			if rec := recover(); rec != nil {
				a.Logger.Error("panic", "path", r.URL.Path, "err", rec)
				WriteErr(w, http.StatusInternalServerError, "internal_error")
			}

			a.Logger.Info(
				"request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		}()

		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type (
	UserKey            struct{}
	HotelIDKey         struct{}
	UserHotelsKey      struct{}
	UserPermissionsKey struct{}
)

func GetUserPermissionsFromContext(ctx context.Context) []models.UserPermissionInfo {
	if perms, ok := ctx.Value(UserPermissionsKey{}).([]models.UserPermissionInfo); ok {
		return perms
	}
	return nil
}

func GetHotelIDFromContext(ctx context.Context) string {
	if hotelID, ok := ctx.Value(HotelIDKey{}).(string); ok {
		return hotelID
	}
	return ""
}

func (a *API) getUserPermissionsFromDB(userID uint, lang string) []models.UserPermissionInfo {
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

func (a *API) GetHotelIDFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie(a.HotelCookie)
	if err != nil || cookie.Value == "" {
		return "", errors.New("hotel not selected")
	}
	return cookie.Value, nil
}

func (a *API) MustGetHotelIDFromCookie(r *http.Request) string {
	hotelID, err := a.GetHotelIDFromCookie(r)
	if err != nil {
		return ""
	}
	return hotelID
}

func (a *API) sessionUser(r *http.Request) (models.SanitizedUser, error) {
	var zero models.SanitizedUser
	cookie, err := r.Cookie(a.SessionCookie)
	if err != nil || cookie.Value == "" {
		return zero, errors.New("missing session")
	}

	var s models.Session
	if err := a.Db.WithContext(r.Context()).Preload("User.UserHotels.Hotel").Preload("User.Roles.Template").Where("id = ? AND expires_at > ?", cookie.Value, time.Now().UTC()).First(&s).Error; err != nil {
		return zero, err
	}
	if !s.User.IsActive {
		return zero, gorm.ErrRecordNotFound
	}
	var roles []models.PermissionTemplate
	for _, v := range s.User.Roles {
		roles = append(roles, v.Template)
	}
	return SanitizeUser(&s.User, roles), nil
}

func SanitizeUser(u *models.User, roles []models.PermissionTemplate) models.SanitizedUser {
	var hotels []models.UserHotelInfo
	for _, uh := range u.UserHotels {
		hotels = append(hotels, models.UserHotelInfo{
			HotelID: uh.HotelID,
			Hotel:   uh.Hotel,
		})
	}
	return models.SanitizedUser{ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName, UserHotels: hotels, Roles: roles}
}
