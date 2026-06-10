package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	spa "hotel"

	"hotel/internal/config"
	"hotel/internal/db"
	"hotel/internal/httpapi"
	system "hotel/internal/httpapi/_system"
	"hotel/internal/httpapi/accounting"
	"hotel/internal/httpapi/admins"
	"hotel/internal/httpapi/auth"
	"hotel/internal/httpapi/common"
	"hotel/internal/httpapi/dashboard"
	"hotel/internal/httpapi/guests"
	"hotel/internal/httpapi/hotels"
	"hotel/internal/httpapi/parking"
	"hotel/internal/httpapi/permissions"
	"hotel/internal/httpapi/reservation"
	"hotel/internal/httpapi/restaurant"
	"hotel/internal/httpapi/rooms"
	sanahttp "hotel/internal/httpapi/sana"
	"hotel/internal/httpapi/services"
	"hotel/internal/httpapi/stays"
	"hotel/internal/httpapi/travelagency"
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

func openAPIHandler(specURL string) http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL(specURL),
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
		"/users":           users.UsersModule{},
		"/rooms":           rooms.RoomsModule{},
		"/guests":          guests.GuestsModule{},
		"/parking":         parking.ParkingModule{},
		"/stays":           stays.StaysModule{},
		"/services":        services.ServicesModule{},
		"/accounting":      accounting.AccountingModule{},
		"/reservation":     reservation.ReservationModule{},
		"/hotels":          hotels.HotelsModule{},
		"/permissions":     permissions.PermissionsModule{},
		"/sana":            sanahttp.New(database, cfg.Sana),
		"/restaurant":      restaurant.RestaurantModule{},
		"/common":          common.CommonModule{},
		"/dashboard":       dashboard.DashboardModule{},
		"/travel-agencies": travelagency.TravelAgencyModule{},
		"/admins":          admins.AdminsModule{},
	})

	spaGroup := fuego.Group(srv, "/")
	spaGroup.Mux.Handle("/", spa.SPAHandler())

	jobsCtx, cancel := context.WithCancel(context.Background())
	go cleanupExpiredSessions(jobsCtx, database, logger)
	go nightAuditJob(jobsCtx, database, logger)

	return &App{cfg: cfg, db: database, server: srv, logger: logger, stopJobs: cancel}, nil
}

func (a *App) Run() error {
	go func() {
		a.logger.Info("server starting", "addr", a.cfg.Addr)
		_ = a.server.Run()
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
			ID:   cfg.SeedHotelCodeName,
			Code: cfg.SeedHotelCodeName,
			Name: cfg.SeedHotelName,
			Address: cfg.SeedHotelAddress,
			Phone:   cfg.SeedHotelPhone,
			Email:   cfg.SeedHotelEmail,
		}).Error; err != nil {
			return fmt.Errorf("create default hotel: %w", err)
		}
	}

	var user models.User
	if err := db.Where("email = ?", cfg.SeedAdminEmail).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("check admin user: %w", err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}

		user = models.User{
			Email:        cfg.SeedAdminEmail,
			PasswordHash: string(hash),
			FirstName:    cfg.SeedAdminFName,
			LastName:     cfg.SeedAdminLName,
			HotelID:      cfg.SeedHotelCodeName,
			IsActive:     true,
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("create admin user: %w", err)
		}
	}

	var uhCount int64
	if err := db.Model(&models.UserHotel{}).Where("user_id = ? AND hotel_id = ?", user.ID, cfg.SeedHotelCodeName).Count(&uhCount).Error; err != nil {
		return fmt.Errorf("check user-hotel link: %w", err)
	}
	if uhCount == 0 {
		if err := db.Create(&models.UserHotel{
			UserID:  user.ID,
			HotelID: cfg.SeedHotelCodeName,
		}).Error; err != nil {
			return fmt.Errorf("create user-hotel link: %w", err)
		}
	}

	var allPermissions []models.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return fmt.Errorf("find all permissions: %w", err)
	}

	if err := db.Where("user_id = ? AND hotel_id = ?", user.ID, cfg.SeedHotelCodeName).Delete(&models.UserPermission{}).Error; err != nil {
		return fmt.Errorf("clear admin permissions: %w", err)
	}

	userPermissions := make([]models.UserPermission, 0, len(allPermissions))
	for _, p := range allPermissions {
		userPermissions = append(userPermissions, models.UserPermission{
			UserID:       user.ID,
			HotelID:      &cfg.SeedHotelCodeName,
			PermissionID: p.ID,
			Granted:      true,
		})
	}
	if err := db.CreateInBatches(userPermissions, len(userPermissions)).Error; err != nil {
		return fmt.Errorf("grant admin permissions: %w", err)
	}

	// Also create an Admin record for the seed admin
	var adminCount int64
	if err := db.Model(&models.Admin{}).Where("email = ?", cfg.SeedAdminEmail).Count(&adminCount).Error; err != nil {
		return fmt.Errorf("check admin record: %w", err)
	}
	if adminCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(cfg.SeedAdminPass), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}
		admin := models.Admin{
			FirstName:    cfg.SeedAdminFName,
			LastName:     cfg.SeedAdminLName,
			Email:        cfg.SeedAdminEmail,
			Username:     cfg.SeedAdminEmail,
			PasswordHash: string(hash),
			IsActive:     true,
			IsSuperAdmin: true,
			AdminHotels: []models.AdminHotel{
				{HotelID: cfg.SeedHotelCodeName},
			},
		}
		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("create admin record: %w", err)
		}
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
			if err := expireUnpaidReservations(db, logger); err != nil {
				logger.Warn("expire reservations failed", "err", err.Error())
			}
		}
	}
}

func expireUnpaidReservations(db *gorm.DB, logger *slog.Logger) error {
	var awaitingPaymentStatus models.ReservationStatus
	if err := db.Where("slug = ?", "awaiting_payment").First(&awaitingPaymentStatus).Error; err != nil {
		return err
	}
	var expiredStatus models.ReservationStatus
	if err := db.Where("slug = ?", "expired").First(&expiredStatus).Error; err != nil {
		return err
	}

	res := db.Model(&models.Reservation{}).
		Where("status_id = ? AND payment_deadline IS NOT NULL AND payment_deadline < ?", awaitingPaymentStatus.ID, time.Now().UTC()).
		Update("status_id", expiredStatus.ID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		logger.Info("expired reservations", "count", res.RowsAffected)
	}
	return nil
}

func nightAuditJob(ctx context.Context, db *gorm.DB, logger *slog.Logger) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now()
			// Only run at night audit hour (default 03:00)
			if now.Hour() != 3 {
				continue
			}

			// Find all hotels and their night audit settings
			var hotels []models.Hotel
			if err := db.Find(&hotels).Error; err != nil {
				logger.Warn("night audit: fetch hotels failed", "err", err.Error())
				continue
			}

			for _, hotel := range hotels {
				var settings models.HotelSetting
				if err := db.Where("hotel_id = ?", hotel.ID).First(&settings).Error; err != nil {
					continue
				}

				// Parse night audit hour
				auditHour := 3
				if settings.NightAuditHour != "" {
					if h, err := time.Parse("15:04", settings.NightAuditHour); err == nil {
						auditHour = h.Hour()
					}
				}

				if now.Hour() != auditHour {
					continue
				}

				// Night audit logic:
				// 1. Mark no-show guests for stays that passed scheduled departure
				var noShowStatus models.StayStatus
				if err := db.Where("slug = ?", "no_show").First(&noShowStatus).Error; err == nil {
					var waitingStatus models.StayStatus
					if err := db.Where("slug = ?", "waiting").First(&waitingStatus).Error; err == nil {
						res := db.Model(&models.Stay{}).
							Where("hotel_id = ? AND status_id = ? AND scheduled_departure_date < ?", hotel.ID, waitingStatus.ID, now.UTC()).
							Update("status_id", noShowStatus.ID)
						if res.Error != nil {
							logger.Warn("night audit: mark no-show failed", "err", res.Error.Error())
						}
						if res.RowsAffected > 0 {
							logger.Info("night audit: marked no-show stays", "hotel", hotel.ID, "count", res.RowsAffected)
						}
					}
				}

				// 2. Set rooms of checked-out stays to cleaning if not already handled
				var checkedOutStatus models.StayStatus
				if err := db.Where("slug = ?", "checked_out").First(&checkedOutStatus).Error; err == nil {
					var cleaningStatus models.RoomStatus
					if err := db.Where("slug = ?", "cleaning").First(&cleaningStatus).Error; err == nil {
						// Find rooms with checked-out stays that haven't been updated yet
						var stays []models.Stay
						if err := db.Where("hotel_id = ? AND status_id = ? AND actual_check_out IS NOT NULL", hotel.ID, checkedOutStatus.ID).
							Find(&stays).Error; err == nil {
							for _, stay := range stays {
								// Check if room is still occupied
								var activeStayCount int64
								db.Model(&models.Stay{}).
									Where("room_id = ? AND status_id IN (SELECT id FROM stay_statuses WHERE slug IN ('waiting', 'resident'))", stay.RoomID).
									Count(&activeStayCount)
								if activeStayCount == 0 {
									db.Model(&models.Room{}).Where("id = ?", stay.RoomID).Update("status_id", cleaningStatus.ID)
								}
							}
						}
					}
				}
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
