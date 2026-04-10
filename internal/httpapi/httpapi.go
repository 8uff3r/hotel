package httpapi

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"hotel/backend/internal/service"
)

type API struct {
	Logger         *slog.Logger
	Db             *gorm.DB
	SessionCookie  string
	RequestTimeout time.Duration
	Services       service.Services
	SessionTTL     time.Duration
}

func RandomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf) + "-" + uuid.NewString(), nil
}
