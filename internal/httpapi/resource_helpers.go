package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"hotel/internal/repository"
	"hotel/internal/service"

	"github.com/go-playground/validator/v10"
)

func (a *API) ListModel(model any, opts *repository.ListOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := newSlicePtr(model)
		if err := a.Services.Crud.List(r.Context(), model, out, opts); err != nil {
			WriteErr(w, 500, "query_failed")
			return
		}
		WriteJSON(w, 200, map[string]any{"data": out})
	}
}

var validate = validator.New()

func BindAndValidate[T any](model T, w http.ResponseWriter, r *http.Request) (T, error) {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		WriteErr(w, http.StatusBadRequest, err.Error())
		return model, err
	}

	// Run validation
	if err := validate.Struct(model); err != nil {
		WriteErr(w, http.StatusBadRequest, err.Error())
		return model, err
	}

	return model, nil
}

func (a *API) CreateModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := BindAndValidate(model, w, r)
		if err != nil {
			slog.Error("", err)
			return
		}
		if err := a.Services.Crud.Create(r.Context(), v); err != nil {
			slog.Error("%s", err)
			WriteErr(w, 400, "create_failed")
			return
		}
		WriteJSON(w, 201, v)
	}
}

func (a *API) GetModel(model any, opts *repository.GetOptions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := ParseID(r.PathValue("id"))
		if err != nil {
			WriteErr(w, 400, "invalid_id")
			return
		}
		entity := newPtr(model)
		if err := a.Services.Crud.GetByID(r.Context(), model, id, entity, opts); err != nil {
			if service.IsNotFound(err) {
				WriteErr(w, 404, "not_found")
				return
			}
			WriteErr(w, 500, "query_failed")
			return
		}
		WriteJSON(w, 200, entity)
	}
}

func (a *API) UpdateModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := ParseID(r.PathValue("id"))
		if err != nil {
			WriteErr(w, 400, "invalid_id")
			return
		}
		updates, err := BindAndValidate(model, w, r)
		if err != nil {
			return
		}
		delete(updates, "id")
		if err := a.Services.Crud.UpdateByID(r.Context(), model, id, updates); err != nil {
			if service.IsNotFound(err) {
				WriteErr(w, 404, "not_found")
				return
			}
			WriteErr(w, 400, "update_failed")
			return
		}
		WriteJSON(w, 200, map[string]bool{"ok": true})
	}
}

func (a *API) DeleteModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := ParseID(r.PathValue("id"))
		if err != nil {
			WriteErr(w, 400, "invalid_id")
			return
		}
		if err := a.Services.Crud.DeleteByID(r.Context(), model, id); err != nil {
			if service.IsNotFound(err) {
				WriteErr(w, 404, "not_found")
				return
			}
			WriteErr(w, 500, "delete_failed")
			return
		}
		WriteJSON(w, 200, map[string]bool{"ok": true})
	}
}
