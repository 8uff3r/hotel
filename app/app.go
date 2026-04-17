package app

import (
	"context"
	"fmt"
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
	"hotel/internal/httpapi/parking"
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
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				JSONFilePath:     "doc/openapi.json",
				SpecURL:          "/swagger/openapi.json",
				SwaggerURL:       "/swagger",
				DisableSwaggerUI: false,
				UIHandler:        openAPIHandler,
			}),
		),
	)

	api := httpapi.API{
		Logger:         logger,
		Db:             database,
		SessionCookie:  cfg.SessionCookie,
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

	system.RegisterRoutes(&api, srv)
	apiGroup := fuego.Group(srv, "/api")
	authGroup := fuego.Group(apiGroup, "/auth")
	auth.RegisterRoutes(&api, authGroup)

	fuego.Use(apiGroup, api.Auth)
	SetupRouter(&api, apiGroup, PathModuleMap{
		"/rooms":       rooms.RoomsModule{},
		"/guests":      guests.GuestsModule{},
		"/users":       users.UsersModule{},
		"/accounting":  accounting.AccountingModule{},
		"/parking":     parking.ParkingModule{},
		"/reservation": reservation.ReservationModule{},
	})

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
	var count int64
	if err := db.Model(&models.User{}).Where("email = ?", cfg.SeedAdminEmail).Count(&count).Error; err != nil {
		return fmt.Errorf("check admin user: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPass), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}
	err = db.Create(&models.User{
		Email:        cfg.SeedAdminEmail,
		PasswordHash: string(hash),
		FirstName:    cfg.SeedAdminFName,
		LastName:     cfg.SeedAdminLName,
		Role:         "admin",
		IsActive:     true,
	}).Error

	return err
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
		subGroup := fuego.Group(group, path, option.Tags(strings.ToUpper(path[1:])))
		mod.RegisterRoutes(api, subGroup)
	}
}
