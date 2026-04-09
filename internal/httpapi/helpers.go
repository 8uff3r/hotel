package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"hotel/backend/internal/models"
)

func newPtr(model any) any {
	switch model.(type) {
	case *models.Room:
		return &models.Room{}
	case *models.Guest:
		return &models.Guest{}
	case *models.Reservation:
		return &models.Reservation{}
	case *models.Account:
		return &models.Account{}
	case *models.Expense:
		return &models.Expense{}
	case *models.Income:
		return &models.Income{}
	case *models.ParkingLot:
		return &models.ParkingLot{}
	case *models.ParkingSpot:
		return &models.ParkingSpot{}
	case *models.Vehicle:
		return &models.Vehicle{}
	case *models.ParkingTransaction:
		return &models.ParkingTransaction{}
	default:
		return nil
	}
}

func newSlicePtr(model any) any {
	switch model.(type) {
	case *models.User:
		return &[]models.User{}
	case *models.Room:
		return &[]models.Room{}
	case *models.Guest:
		return &[]models.Guest{}
	case *models.Reservation:
		return &[]models.Reservation{}
	case *models.Account:
		return &[]models.Account{}
	case *models.Expense:
		return &[]models.Expense{}
	case *models.Income:
		return &[]models.Income{}
	case *models.ParkingLot:
		return &[]models.ParkingLot{}
	case *models.ParkingSpot:
		return &[]models.ParkingSpot{}
	case *models.Vehicle:
		return &[]models.Vehicle{}
	case *models.ParkingTransaction:
		return &[]models.ParkingTransaction{}
	default:
		return &[]map[string]any{}
	}
}

func normalizeUpdates(in map[string]any) {
	for k, v := range in {
		nk := k
		if !strings.ContainsRune(k, '_') {
			nk = camelToSnake(k)
		}
		if nk != k {
			delete(in, k)
			in[nk] = v
		}
	}
}

func camelToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func parseID(raw string) (uint, error) {
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func decode[T any](r *http.Request, dst T, w http.ResponseWriter) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	// dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid_payload")
		return err
	}
	return nil
}

func writeErr(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func cookieValue(c *http.Cookie) string {
	if c == nil {
		return ""
	}
	return c.Value
}

func defaultStr(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func pickString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok {
			return v
		}
	}
	return ""
}
