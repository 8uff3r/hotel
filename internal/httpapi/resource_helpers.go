package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (a *API) ListModel(model any, preload []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out := newSlicePtr(model)
		db := a.Db.WithContext(r.Context()).Model(model)
		for _, v := range preload {
			db = db.Preload(v)
		}
		if err := db.Order("id DESC").Find(out).Error; err != nil {
			WriteErr(w, 500, "query_failed")
			return
		}
		WriteJSON(w, 200, map[string]any{"data": out})
	}
}

var validate = validator.New()

func BindAndValidate[T any](model *T, w http.ResponseWriter, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		WriteErr(w, http.StatusBadRequest, err.Error())
		return err
	}
	if err := validate.Struct(model); err != nil {
		WriteErr(w, http.StatusBadRequest, err.Error())
		return err
	}
	return nil
}

func (a *API) CreateModel(model any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := BindAndValidate(&model, w, r)
		if err != nil {
			return
		}
		if err := a.Db.WithContext(r.Context()).Create(model).Error; err != nil {
			WriteErr(w, 400, "create_failed")
			return
		}
		WriteJSON(w, 201, model)
	}
}

func (a *API) GetModel(model any, preload []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := ParseID(r.PathValue("id"))
		if err != nil {
			WriteErr(w, 400, "invalid_id")
			return
		}
		entity := newPtr(model)
		db := a.Db.WithContext(r.Context()).Model(model)
		for _, v := range preload {
			db = db.Preload(v)
		}
		if err := db.First(entity, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
		err = BindAndValidate(&model, w, r)
		if err != nil {
			return
		}
		modelValue := reflect.ValueOf(model).Elem()
		modelValue.FieldByName("ID").SetUint(uint64(id))

		res := a.Db.WithContext(r.Context()).Model(model).Where("id = ?", id).Updates(model)
		if res.Error != nil {
			WriteErr(w, 400, "update_failed")
			return
		}
		if res.RowsAffected == 0 {
			WriteErr(w, 404, "not_found")
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
		res := a.Db.WithContext(r.Context()).Delete(model, id)
		if res.Error != nil {
			WriteErr(w, 500, "delete_failed")
			return
		}
		if res.RowsAffected == 0 {
			WriteErr(w, 404, "not_found")
			return
		}
		WriteJSON(w, 200, map[string]bool{"ok": true})
	}
}
