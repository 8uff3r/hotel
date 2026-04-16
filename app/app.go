package app

import (
	"context"
	"errors"
	"fmt"
	spa "hotel"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hotel/internal/config"
	"hotel/internal/db"
	"hotel/internal/httpapi"
	"hotel/internal/models"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type App struct {
	cfg      config.Config
	db       *gorm.DB
	server   *http.Server
	logger   *slog.Logger
	stopJobs context.CancelFunc
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

	r := chi.NewRouter()
	api := httpapi.API{
		Logger:         logger,
		SessionCookie:  cfg.SessionCookie,
		RequestTimeout: cfg.RequestTimeout,
		Db:             database,
		SessionTTL:     cfg.SessionTTL,
	}
	handler := NewRouter(&api, r)

	notFoundHandlerFunc := spa.SPAHandler()
	r.NotFound(notFoundHandlerFunc.ServeHTTP)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	jobsCtx, cancel := context.WithCancel(context.Background())
	go cleanupExpiredSessions(jobsCtx, database, logger)

	return &App{cfg: cfg, db: database, server: server, logger: logger, stopJobs: cancel}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)
	go func() {
		a.logger.Info("server starting", "addr", a.cfg.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-sigCtx.Done():
		a.logger.Info("shutdown signal received")
	case err := <-errCh:
		return fmt.Errorf("server failure: %w", err)
	}

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
