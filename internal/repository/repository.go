package repository

import (
	"context"
	"time"

	"hotel/internal/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserBySession(ctx context.Context, sessionID string) (*models.User, error)
	CreateSession(ctx context.Context, session models.Session) error
	DeleteSession(ctx context.Context, sessionID string) error
	CleanupExpired(ctx context.Context) error
}

type CrudRepository interface {
	List(ctx context.Context, model any, out any, opts *ListOptions) error
	GetByID(ctx context.Context, model any, id uint, out any, opts *GetOptions) error
	Create(ctx context.Context, model any) error
	UpdateByID(ctx context.Context, model any, id uint, updates map[string]any) error
	DeleteByID(ctx context.Context, model any, id uint) error
}

type ReservationRepository interface {
	CheckIn(ctx context.Context, id uint) error
	CheckOut(ctx context.Context, id uint) error
}

type ParkingTxRepository interface {
	CheckOut(ctx context.Context, id uint) error
}

type ParkingRepository interface {
	Stats(ctx context.Context) (*models.ParkingStats, error)
}

type Repositories struct {
	Auth        AuthRepository
	Crud        CrudRepository
	Reservation ReservationRepository
	ParkingTx   ParkingTxRepository
	Parking     ParkingRepository
}

func New(db *gorm.DB) Repositories {
	return Repositories{
		Auth:        &authRepo{db: db},
		Crud:        &crudRepo{db: db},
		Reservation: &reservationRepo{db: db},
		ParkingTx:   &parkingTxRepo{db: db},
		Parking:     &parkingRepo{db: db},
	}
}

type authRepo struct{ db *gorm.DB }

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) GetUserBySession(ctx context.Context, sessionID string) (*models.User, error) {
	var s models.Session
	if err := r.db.WithContext(ctx).Preload("User").Where("id = ? AND expires_at > ?", sessionID, time.Now().UTC()).First(&s).Error; err != nil {
		return nil, err
	}
	if !s.User.IsActive {
		return nil, gorm.ErrRecordNotFound
	}
	return &s.User, nil
}

func (r *authRepo) CreateSession(ctx context.Context, session models.Session) error {
	return r.db.WithContext(ctx).Create(&session).Error
}
func (r *authRepo) DeleteSession(ctx context.Context, sessionID string) error {
	return r.db.WithContext(ctx).Delete(&models.Session{}, "id = ?", sessionID).Error
}
func (r *authRepo) CleanupExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Delete(&models.Session{}, "expires_at <= ?", time.Now().UTC()).Error
}

type crudRepo struct{ db *gorm.DB }

type ListOptions struct {
	Preload []string
}

func (r *crudRepo) List(ctx context.Context, model any, out any, opts *ListOptions) error {
	db := r.db.WithContext(ctx).Model(model)
	if opts != nil && opts.Preload != nil {
		for _, v := range opts.Preload {
			db = db.Preload(v)
		}
	}

	return db.Order("id DESC").Find(out).Error

}

type GetOptions struct {
	Preload []string
}

func (r *crudRepo) GetByID(ctx context.Context, model any, id uint, out any, opts *GetOptions) error {
	db := r.db.WithContext(ctx).Model(model)
	if opts != nil && opts.Preload != nil {
		for _, v := range opts.Preload {
			db = db.Preload(v)
		}
	}

	return db.First(out, id).Error
}
func (r *crudRepo) Create(ctx context.Context, model any) error {
	return r.db.WithContext(ctx).Create(model).Error
}
func (r *crudRepo) UpdateByID(ctx context.Context, model any, id uint, updates map[string]any) error {
	res := r.db.WithContext(ctx).Model(model).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *crudRepo) DeleteByID(ctx context.Context, model any, id uint) error {
	res := r.db.WithContext(ctx).Delete(model, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type reservationRepo struct{ db *gorm.DB }

func (r *reservationRepo) CheckIn(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_in", "actual_check_in": time.Now().UTC()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *reservationRepo) CheckOut(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_out", "actual_check_out": time.Now().UTC()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type parkingTxRepo struct{ db *gorm.DB }

func (r *parkingTxRepo) CheckOut(ctx context.Context, id uint) error {
	now := time.Now().UTC()
	res := r.db.WithContext(ctx).Model(&models.ParkingTransaction{}).Where("id = ?", id).Updates(map[string]any{"status": "completed", "exit_time": now})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

type parkingRepo struct{ db *gorm.DB }

func (r *parkingRepo) Stats(ctx context.Context) (*models.ParkingStats, error) {
	db := r.db.WithContext(ctx)
	var totalLots int64
	if err := db.Model(&models.ParkingLot{}).Count(&totalLots).Error; err != nil {
		return nil, err
	}

	var totalSpots int64
	if err := db.Model(&models.ParkingSpot{}).Count(&totalSpots).Error; err != nil {
		return nil, err
	}

	var availableSpots int64
	if err := db.Model(&models.ParkingSpot{}).Where("status = ?", "available").Count(&availableSpots).Error; err != nil {
		return nil, err
	}

	println(totalLots, totalSpots, availableSpots)

	return &models.ParkingStats{Lots: totalLots, Spots: totalSpots, AvailableSpots: availableSpots}, nil
}
