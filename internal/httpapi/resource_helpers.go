package httpapi

import (
	"encoding/json"
	"errors"
	"hotel/internal/models"
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

type listConfig struct {
	preload         []string
	translate       bool
	translateFields []string
}
type ListOption func(*listConfig)

func WithPreload(preload ...string) ListOption {
	return func(c *listConfig) {
		c.preload = preload
	}
}

func WithFieldTranslation(fields ...string) ListOption {
	return func(lc *listConfig) {
		lc.translateFields = append(lc.translateFields, fields...)
	}
}

func WithTranslation() ListOption {
	return func(lc *listConfig) {
		lc.translate = true
	}
}

func ListModel[T any](db *gorm.DB, model T, opts ...ListOption) FuegoHandler[PaginatedResponse[T], any, Params] {
	lc := listConfig{}
	for _, v := range opts {
		v(&lc)
	}
	return listModel(db, model, lc.preload, lc.translate, lc.translateFields)
}
func listModel[T any](db *gorm.DB, model T, preload []string, translate bool, translateFields []string) FuegoHandler[PaginatedResponse[T], any, Params] {
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

		lang := c.Header("Accept-Language")
		if lang == "" {
			lang = "fa"
		}
		if translate {
			applyTranslations(&out, lang)
			applyFieldTranslations(&out, lang)
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
func applyTranslationsReflection[T any](items *[]T, lang string) {
	for i := range *items {
		item := &(*items)[i]
		if translatable, ok := any(item).(models.Translatable); ok {
			applyTranslationOnTranslatable(translatable, lang)
		}
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

func GetModel[T any](db *gorm.DB, model T, opts ...ListOption) FuegoAnyHandler[T] {
	lc := listConfig{}
	for _, v := range opts {
		v(&lc)
	}
	return getModel(db, model, lc.preload, lc.translate)
}

func getModel[T any](db *gorm.DB, model T, preload []string, translate bool) FuegoAnyHandler[T] {
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
		if err := db.First(&entity, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zero, fuego.NotFoundError{}
			}
			return zero, nil
		}
		lang := c.Header("Accept-Language")
		if lang == "" {
			lang = "fa"
		}
		slice := []T{entity}
		if translate {
			applyTranslations(&slice, lang)
			applyFieldTranslations(&slice, lang)
		}
		return slice[0], nil
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
		rv := reflect.ValueOf(body)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem() // get the struct from the pointer
		} else {
			rv = reflect.ValueOf(&body).Elem() // addressable struct from value
		}
		idField := rv.FieldByName("ID")
		if !idField.IsValid() {
			return zero, fuego.BadRequestError{Title: "model has no ID field"}
		}
		idField.SetUint(uint64(id))

		res := db.WithContext(c).Model(model).Where("id = ?", id).Updates(body)
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
func applyTranslationOnTranslatable(t models.Translatable, lang string) {
	translations := t.GetTranslation()
	if lang != "en" {
		if translated, exists := translations[lang]; exists {
			t.SetName(translated)
		}
	}
	t.ClearTranslation()
}

func applyTranslations[T any](items *[]T, lang string) {
	for i := range *items {
		item := &(*items)[i]
		if translatable, ok := any(item).(models.Translatable); ok {
			applyTranslationOnTranslatable(translatable, lang)
		}
	}
}

func applyFieldTranslations[T any](items *[]T, lang string) {
	for i := range *items {
		item := &(*items)[i]

		if container, ok := any(item).(models.HasTranslatables); ok {
			translatables := container.GetTranslatables()

			for _, tr := range translatables {
				applyTranslationOnTranslatable(tr, lang)
			}
		}
	}
}
