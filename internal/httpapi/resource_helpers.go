package httpapi

import (
	"log/slog"
	"net/http"

	"hotel/backend/internal/service"
)

func (a *API) listModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := newSlicePtr(model)
		if err := a.services.Crud.List(r.Context(), model, out); err != nil {
			writeErr(w, 500, "query_failed")
			return
		}
		writeJSON(w, 200, map[string]any{"data": out})
	}
}

func (a *API) createModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entity := newPtr(model)
		if err := decode(r, entity, w); err != nil {
			slog.Error("", err)
			return
		}
		if err := a.services.Crud.Create(r.Context(), entity); err != nil {
			slog.Error("%s", err)
			writeErr(w, 400, "create_failed")
			return
		}
		writeJSON(w, 201, entity)
	}
}

func (a *API) getModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.PathValue("id"))
		if err != nil {
			writeErr(w, 400, "invalid_id")
			return
		}
		entity := newPtr(model)
		if err := a.services.Crud.GetByID(r.Context(), model, id, entity); err != nil {
			if service.IsNotFound(err) {
				writeErr(w, 404, "not_found")
				return
			}
			writeErr(w, 500, "query_failed")
			return
		}
		writeJSON(w, 200, entity)
	}
}

func (a *API) updateModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.PathValue("id"))
		if err != nil {
			writeErr(w, 400, "invalid_id")
			return
		}
		updates := map[string]any{}
		if err := decode(r, &updates, w); err != nil {
			return
		}
		delete(updates, "id")
		normalizeUpdates(updates)
		if err := a.services.Crud.UpdateByID(r.Context(), model, id, updates); err != nil {
			if service.IsNotFound(err) {
				writeErr(w, 404, "not_found")
				return
			}
			writeErr(w, 400, "update_failed")
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	}
}

func (a *API) deleteModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseID(r.PathValue("id"))
		if err != nil {
			writeErr(w, 400, "invalid_id")
			return
		}
		if err := a.services.Crud.DeleteByID(r.Context(), model, id); err != nil {
			if service.IsNotFound(err) {
				writeErr(w, 404, "not_found")
				return
			}
			writeErr(w, 500, "delete_failed")
			return
		}
		writeJSON(w, 200, map[string]bool{"ok": true})
	}
}
