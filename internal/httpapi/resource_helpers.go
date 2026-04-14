package httpapi

import (
	"log/slog"
	"net/http"

	"hotel/internal/repository"
	"hotel/internal/service"
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

func (a *API) CreateModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity := newPtr(model)
		if err := Decode(r, entity, w); err != nil {
			slog.Error("", err)
			return
		}
		if err := a.Services.Crud.Create(r.Context(), entity); err != nil {
			slog.Error("%s", err)
			WriteErr(w, 400, "create_failed")
			return
		}
		WriteJSON(w, 201, entity)
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
		updates := map[string]any{}
		if err := Decode(r, &updates, w); err != nil {
			return
		}
		delete(updates, "id")
		normalizeUpdates(updates)
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
