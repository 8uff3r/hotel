package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type FuegoAnyHandler[T any] = func(fuego.Context[any, any]) (T, error)
type FuegoHandler[T any, B any, P any] = func(fuego.Context[B, P]) (T, error)
type Params struct {
	Limit int `query:"limit"`
	page  int `query:"page"`
}

func ListModel[T any](db *gorm.DB, model T, preload []string) FuegoHandler[[]T, any, Params] {
	return func(c fuego.ContextWithParams[Params]) ([]T, error) {
		out := []T{}
		db = db.WithContext(c).Model(model)
		for _, v := range preload {
			db = db.Preload(v)
		}
		if err := db.Order("id DESC").Find(out).Error; err != nil {
			return nil, err
		}
		return out, nil
	}
}

func CreateModel[T any](db *gorm.DB, model T) FuegoHandler[T, T, any] {
	return func(c fuego.ContextWithBody[T]) (T, error) {
		var zero T
		body, err := c.Body()
		if err != nil {
			return zero, err
		}
		if err := db.WithContext(c).Create(body).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		return model, nil
	}
}

func GetModel[T any](db *gorm.DB, model T, preload []string) FuegoAnyHandler[T] {
	return func(c fuego.ContextNoBody) (T, error) {
		var zero T
		id, err := ParseID(c.PathParam("id"))
		if err != nil {
			return zero, err
		}
		var entity T
		db := db.WithContext(c).Model(model)
		for _, v := range preload {
			db = db.Preload(v)
		}
		if err := db.First(entity, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zero, fuego.NotFoundError{}
			}
			return zero, nil
		}
		return entity, nil
	}
}

func UpdateModel[T any](db *gorm.DB, model T) FuegoHandler[T, T, any] {
	return func(c fuego.ContextWithBody[T]) (T, error) {
		var zero T
		id, err := ParseID(c.PathParam("id"))
		if err != nil {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}

		body, err := c.Body()
		if err != nil {
			return zero, err
		}
		modelValue := reflect.ValueOf(body).Elem()
		modelValue.FieldByName("ID").SetUint(uint64(id))

		res := db.WithContext(c).Model(model).Where("id = ?", id).Updates(model)
		if res.Error != nil {
			return zero, fuego.BadRequestError{Title: "update_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return body, nil
	}
}

type okResponse struct{ ok bool }

type deleteDto struct{ id string }

func DeleteModel(db *gorm.DB, model any) FuegoHandler[okResponse, any, deleteDto] {
	var zero okResponse
	return func(c fuego.ContextWithParams[deleteDto]) (okResponse, error) {
		id, err := ParseID(c.PathParam("id"))
		if err != nil {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		res := db.WithContext(c).Delete(model, id)
		if res.Error != nil {
			return zero, fuego.InternalServerError{Title: "delete_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return okResponse{ok: true}, nil
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
