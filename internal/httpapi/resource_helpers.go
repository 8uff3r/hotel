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
	Page  int `query:"page"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func ListModel[T any](db *gorm.DB, model T, preload []string) FuegoHandler[PaginatedResponse[T], any, Params] {
	return func(c fuego.ContextWithParams[Params]) (PaginatedResponse[T], error) {
		page := c.QueryParamInt("page")
		if page < 1 {
			page = 1
		}
		limit := c.QueryParamInt("limit")
		if limit < 1 || limit > 100 {
			limit = 20
		}
		offset := (page - 1) * limit

		out := []T{}
		q := db.WithContext(c).Model(model)
		for _, v := range preload {
			q = q.Preload(v)
		}

		var total int64
		if err := q.Count(&total).Error; err != nil {
			return PaginatedResponse[T]{}, err
		}

		if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
			return PaginatedResponse[T]{}, err
		}

		totalPages := int((total + int64(limit) - 1) / int64(limit))
		return PaginatedResponse[T]{
			Data:       out,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		}, nil
	}
}

func CreateModel[T any](db *gorm.DB, model T) FuegoHandler[T, T, any] {
	return func(c fuego.ContextWithBody[T]) (T, error) {
		var zero T
		body, err := c.Body()
		if err != nil {
			return zero, err
		}
		if err := db.WithContext(c.Context()).Create(&body).Error; err != nil {
			return zero, fuego.BadRequestError{Title: "create_failed"}
		}
		return body, nil
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
