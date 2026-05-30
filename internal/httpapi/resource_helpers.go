package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"hotel/internal/models"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type (
	FuegoAnyHandler[T any]            = func(fuego.Context[any, any]) (T, error)
	FuegoHandler[T any, B any, P any] = func(fuego.Context[B, P]) (T, error)
	Params                            struct {
		Limit   int    `query:"limit"`
		Page    int    `query:"page"`
		Filters string `query:"filters"`
	}
)

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}
type Filter struct {
	Query any
	Args  []any
}

// QueryOption modifies the gorm query
type QueryOption func(fuego.ContextWithParams[Params], *gorm.DB) *gorm.DB

// PostProcessOption runs after the query on the results
type PostProcessOption[T any] func(fuego.ContextWithParams[Params], *[]T)

func WithPreload(relations ...string) QueryOption {
	return func(_ctx fuego.ContextWithParams[Params], db *gorm.DB) *gorm.DB {
		for _, r := range relations {
			db = db.Preload(r)
		}
		return db
	}
}

func WithAllowedFilters(allowed ...string) QueryOption {
	return func(ctx fuego.ContextWithParams[Params], db *gorm.DB) *gorm.DB {
		raw := ctx.PathParam("filters")
		return applyQueryFilters(db, raw, allowed)
	}
}

// func WithHotelScope(fn func(context.Context) uint) QueryOption {
// 	return func(ctx context.Context, db *gorm.DB) *gorm.DB {
// 		if id := fn(ctx); id > 0 {
// 			return db.Where("hotel_id = ?", id)
// 		}
// 		return db
// 	}
// }

func WithTranslation[T any]() PostProcessOption[T] {
	return func(ctx fuego.ContextWithParams[Params], out *[]T) {
		lang := ctx.Header("Accept-Language")
		if lang == "" {
			lang = "fa"
		}
		models.ApplyTranslations(out, lang)
		models.ApplyFieldTranslations(out, lang)
	}
}

func ListModel[T any](db *gorm.DB, opts ...any) FuegoHandler[PaginatedResponse[T], any, Params] {
	return func(c fuego.ContextWithParams[Params]) (PaginatedResponse[T], error) {
		q := db.WithContext(c).Model(new(T))

		var postOpts []PostProcessOption[T]
		for _, opt := range opts {
			switch o := opt.(type) {
			case QueryOption:
				q = o(c, q)
			case PostProcessOption[T]:
				postOpts = append(postOpts, o)
			}
		}

		result, err := Paginate[T](c, q)
		if err != nil {
			return result, err
		}

		for _, p := range postOpts {
			p(c, &result.Data)
		}

		return result, nil
	}
}

func Paginate[T any](c fuego.ContextWithParams[Params], q *gorm.DB) (PaginatedResponse[T], error) {
	page := max(c.QueryParamInt("page"), 1)
	limit, err := c.QueryParamIntErr("limit")
	if limit == 0 || err != nil {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	var out []T
	if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&out).Error; err != nil {
		return PaginatedResponse[T]{}, err
	}

	return PaginatedResponse[T]{
		Data:       out,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}, nil
}

func CreateModel[T any](db *gorm.DB) FuegoHandler[T, T, any] {
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

type GetOption func(*gorm.DB) *gorm.DB

func GetModel[T any](db *gorm.DB, preloads ...string) FuegoAnyHandler[T] {
	return func(c fuego.ContextNoBody) (T, error) {
		var zero T
		id, err := ParseID(c.PathParam("id"))
		if err != nil {
			return zero, err
		}

		q := db.WithContext(c).Model(new(T))

		var entity T

		for _, r := range preloads {
			db = db.Preload(r)
		}
		if err := q.First(&entity, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return zero, fuego.NotFoundError{}
			}
			return zero, err
		}

		out := []T{entity}

		lang := c.Header("Accept-Language")
		if lang == "" {
			lang = "fa"
		}
		models.ApplyTranslations(&out, lang)
		models.ApplyFieldTranslations(&out, lang)
		return out[0], nil
	}
}

func UpdateModel[T any](db *gorm.DB) FuegoHandler[T, T, any] {
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
		if rv.Kind() == reflect.Pointer {
			rv = rv.Elem() // get the struct from the pointer
		} else {
			rv = reflect.ValueOf(&body).Elem() // addressable struct from value
		}
		idField := rv.FieldByName("ID")
		if !idField.IsValid() {
			return zero, fuego.BadRequestError{Title: "model has no ID field"}
		}
		idField.SetUint(uint64(id))

		res := db.WithContext(c).Model(new(T)).Where("id = ?", id).Updates(body)
		if res.Error != nil {
			return zero, fuego.BadRequestError{Title: "update_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return body, nil
	}
}

type okResponse struct {
	Ok bool `json:"ok"`
}

type deleteDto struct{ id string }

func DeleteModel[T any](db *gorm.DB) FuegoHandler[okResponse, any, deleteDto] {
	var zero okResponse
	return func(c fuego.ContextWithParams[deleteDto]) (okResponse, error) {
		id, err := ParseID(c.PathParam("id"))
		if err != nil {
			return zero, fuego.BadRequestError{Title: "invalid_id"}
		}
		res := db.WithContext(c).Delete(new(T), id)
		if res.Error != nil {
			return zero, fuego.InternalServerError{Title: "delete_failed"}
		}
		if res.RowsAffected == 0 {
			return zero, fuego.NotFoundError{}
		}
		return okResponse{Ok: true}, nil
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

func applyQueryFilters(db *gorm.DB, raw string, allowed []string) *gorm.DB {
	allowedSet := make(map[string]struct{}, len(allowed))
	for _, col := range allowed {
		allowedSet[col] = struct{}{}
	}

	ops := map[string]string{
		"eq":   "= ?",
		"gt":   "> ?",
		"lt":   "< ?",
		"gte":  ">= ?",
		"lte":  "<= ?",
		"like": "LIKE ?",
		"neq":  "!= ?",
	}

	for segment := range strings.SplitSeq(raw, ";") {
		parts := strings.SplitN(segment, ":", 3)
		if len(parts) != 3 {
			continue
		}
		col, op, val := parts[0], parts[1], parts[2]

		if _, ok := allowedSet[col]; !ok {
			continue // silently skip disallowed columns
		}
		clause, ok := ops[op]
		if !ok {
			continue
		}
		if op == "like" {
			val = "%" + val + "%"
		}
		db = db.Where(col+" "+clause, val)
	}
	return db
}
