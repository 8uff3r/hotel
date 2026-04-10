package users

import (
	h "hotel/backend/internal/httpapi"
	"hotel/backend/internal/models"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type UsersModule struct {
	*h.API
}

func (m UsersModule) RegisterRoutes(api *h.API, r chi.Router) {
	u := UsersModule{api}

	r.Get("/", u.usersList)
	r.Post("/", u.usersCreate)
}

func (u *UsersModule) usersList(w http.ResponseWriter, r *http.Request) {
	var rows []models.User
	if err := u.Services.Crud.List(r.Context(), &models.User{}, &rows, nil); err != nil {
		h.WriteErr(w, 500, "query_failed")
		return
	}
	out := make([]any, 0, len(rows))
	for i := range rows {
		out = append(out, h.SanitizeUser(&rows[i]))
	}
	h.WriteJSON(w, 200, map[string]any{"data": out})
}

func (u *UsersModule) usersCreate(w http.ResponseWriter, r *http.Request) {
	var in map[string]any
	if err := h.Decode(r, &in, w); err != nil {
		return
	}
	email, _ := in["email"].(string)
	password, _ := in["password"].(string)
	firstName := h.PickString(in, "firstName", "first_name")
	lastName := h.PickString(in, "lastName", "last_name")
	role := h.PickString(in, "role")
	if role == "" {
		if roles, ok := in["roles"].([]any); ok && len(roles) > 0 {
			if r0, ok := roles[0].(string); ok {
				role = r0
			}
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.WriteErr(w, 500, "hash_failed")
		return
	}
	user := &models.User{
		Email:        strings.TrimSpace(email),
		PasswordHash: string(hash),
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		Role:         h.DefaultStr(role, "staff"),
		IsActive:     true,
	}
	if v, ok := in["isActive"].(bool); ok {
		user.IsActive = v
	}
	if err := u.Services.Crud.Create(r.Context(), user); err != nil {
		h.WriteErr(w, 400, "create_failed")
		return
	}
	h.WriteJSON(w, 201, map[string]any{"id": user.ID})
}
