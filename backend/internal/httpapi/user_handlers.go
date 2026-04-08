package httpapi

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"hotel/backend/internal/models"
)

func (a *API) usersList(w http.ResponseWriter, r *http.Request) {
	var rows []models.User
	if err := a.services.Crud.List(r.Context(), &models.User{}, &rows); err != nil {
		writeErr(w, 500, "query_failed")
		return
	}
	out := make([]any, 0, len(rows))
	for i := range rows {
		out = append(out, sanitizeUser(&rows[i]))
	}
	writeJSON(w, 200, map[string]any{"data": out})
}

func (a *API) usersCreate(w http.ResponseWriter, r *http.Request) {
	var in map[string]any
	if !decode(r, &in, w) {
		return
	}
	email, _ := in["email"].(string)
	password, _ := in["password"].(string)
	firstName := pickString(in, "firstName", "first_name")
	lastName := pickString(in, "lastName", "last_name")
	role := pickString(in, "role")
	if role == "" {
		if roles, ok := in["roles"].([]any); ok && len(roles) > 0 {
			if r0, ok := roles[0].(string); ok {
				role = r0
			}
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		writeErr(w, 500, "hash_failed")
		return
	}
	user := &models.User{
		Email:        strings.TrimSpace(email),
		PasswordHash: string(hash),
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		Role:         defaultStr(role, "staff"),
		IsActive:     true,
	}
	if v, ok := in["isActive"].(bool); ok {
		user.IsActive = v
	}
	if err := a.services.Crud.Create(r.Context(), user); err != nil {
		writeErr(w, 400, "create_failed")
		return
	}
	writeJSON(w, 201, map[string]any{"id": user.ID})
}
