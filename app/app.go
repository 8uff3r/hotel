package app

import (
	"context"
	"fmt"
	spa "hotel"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"hotel/internal/config"
	"hotel/internal/db"
	"hotel/internal/httpapi"
	system "hotel/internal/httpapi/_system"
	"hotel/internal/httpapi/accounting"
	"hotel/internal/httpapi/auth"
	"hotel/internal/httpapi/guests"
	"hotel/internal/httpapi/hotels"
	"hotel/internal/httpapi/parking"
	"hotel/internal/httpapi/permissions"
	"hotel/internal/httpapi/reservation"
	"hotel/internal/httpapi/rooms"
	"hotel/internal/httpapi/users"
	"hotel/internal/models"

	fuego "github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type App struct {
	cfg      config.Config
	db       *gorm.DB
	server   *fuego.Server
	logger   *slog.Logger
	stopJobs context.CancelFunc
}

func openAPIHandler(specUrl string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL(specUrl),
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.PersistAuthorization(true),
	)
}

func New(cfg config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	database, err := db.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	if err := ensureAdmin(database, cfg); err != nil {
		return nil, err
	}

	srv := fuego.NewServer(
		fuego.WithAddr(cfg.Addr),
		fuego.WithoutAutoGroupTags(),
		fuego.WithoutStartupMessages(),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				JSONFilePath:     "doc/openapi.json",
				SpecURL:          "/swagger/openapi.json",
				SwaggerURL:       "/swagger",
				DisableSwaggerUI: false,
				UIHandler:        openAPIHandler,
				PrettyFormatJSON: true,
			}),
		),
	)

	api := httpapi.API{
		Logger:         logger,
		Db:             database,
		SessionCookie:  cfg.SessionCookie,
		HotelCookie:    cfg.HotelCookie,
		RequestTimeout: cfg.RequestTimeout,
		SessionTTL:     cfg.SessionTTL,
	}

	// middleware := func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		start := time.Now()

	// 		defer func() {
	// 			if rec := recover(); rec != nil {
	// 				logger.Error("panic", "path", r.URL.Path, "err", rec)
	// 				http.Error(w, "internal_error", http.StatusInternalServerError)
	// 			}

	// 			logger.Info("request",
	// 				"method", r.Method,
	// 				"path", r.URL.Path,
	// 				"duration_ms", time.Since(start).Milliseconds(),
	// 			)
	// 		}()

	// 		w.Header().Set("Content-Type", "application/json")
	// 		next.ServeHTTP(w, r)
	// 	})
	// }
	// fuego.Use(srv, middleware)

	apiGroup := fuego.Group(srv, "/api")
	SetupRouter(&api, apiGroup, PathModuleMap{
		"/":     system.SystemModule{},
		"/auth": auth.AuthModule{},
	})

	fuego.Use(apiGroup, api.AuthMiddleware)
	SetupRouter(&api, apiGroup, PathModuleMap{
		"/users":       users.UsersModule{},
		"/rooms":       rooms.RoomsModule{},
		"/guests":      guests.GuestsModule{},
		"/parking":     parking.ParkingModule{},
		"/accounting":  accounting.AccountingModule{},
		"/reservation": reservation.ReservationModule{},
		"/hotels":      hotels.HotelsModule{},
		"/permissions": permissions.PermissionsModule{},
	})

	spaGroup := fuego.Group(srv, "/")
	spaGroup.Mux.Handle("/", spa.SPAHandler())

	jobsCtx, cancel := context.WithCancel(context.Background())
	go cleanupExpiredSessions(jobsCtx, database, logger)

	return &App{cfg: cfg, db: database, server: srv, logger: logger, stopJobs: cancel}, nil
}

func (a *App) Run() error {
	go func() {
		a.logger.Info("server starting", "addr", a.cfg.Addr)
		a.server.Run()
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-sigCtx.Done()
	a.logger.Info("shutdown signal received")

	a.stopJobs()
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	sqlDB, err := a.db.DB()
	if err != nil {
		return fmt.Errorf("sql db: %w", err)
	}
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("db close: %w", err)
	}

	a.logger.Info("server stopped")
	return nil
}

func ensureAdmin(db *gorm.DB, cfg config.Config) error {
	var hotelCount int64
	if err := db.Model(&models.Hotel{}).Count(&hotelCount).Error; err != nil {
		return fmt.Errorf("check hotel: %w", err)
	}
	if hotelCount == 0 {
		if err := db.Create(&models.Hotel{
			ID:      cfg.SeedHotelCodeName,
			Name:    cfg.SeedHotelName,
			Address: cfg.SeedHotelAddress,
			Phone:   cfg.SeedHotelPhone,
			Email:   cfg.SeedHotelEmail,
		}).Error; err != nil {
			return fmt.Errorf("create default hotel: %w", err)
		}
	}

	var userCount int64
	if err := db.Model(&models.User{}).Where("email = ?", cfg.SeedAdminEmail).Count(&userCount).Error; err != nil {
		return fmt.Errorf("check admin user: %w", err)
	}
	if userCount > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	user := &models.User{
		Email:        cfg.SeedAdminEmail,
		PasswordHash: string(hash),
		FirstName:    cfg.SeedAdminFName,
		LastName:     cfg.SeedAdminLName,
		IsActive:     true,
	}
	if err := db.Create(user).Error; err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	if err := db.Create(&models.UserHotel{
		UserID:  user.ID,
		HotelID: cfg.SeedHotelCodeName,
	}).Error; err != nil {
		return fmt.Errorf("create user-hotel link: %w", err)
	}

	return nil
}

func cleanupExpiredSessions(ctx context.Context, db *gorm.DB, logger *slog.Logger) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := db.WithContext(ctx).Delete(&models.Session{}, "expires_at <= ?", time.Now().UTC()).Error; err != nil {
				logger.Warn("cleanup sessions failed", "err", err.Error())
			}
		}
	}
}

type ModuleRouter interface {
	RegisterRoutes(*httpapi.API, *fuego.Server)
}

type PathModuleMap map[string]ModuleRouter

func SetupRouter(api *httpapi.API, group *fuego.Server, modules PathModuleMap) {
	for path, mod := range modules {
		subGroup := fuego.Group(group, path, option.Tags(strings.ToUpper(path)))
		mod.RegisterRoutes(api, subGroup)
	}
}
