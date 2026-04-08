package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"hotel/backend/internal/models"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserBySession(ctx context.Context, sessionID string) (*models.User, error)
	CreateSession(ctx context.Context, session models.Session) error
	DeleteSession(ctx context.Context, sessionID string) error
	CleanupExpired(ctx context.Context) error
}

type CrudRepository interface {
	List(ctx context.Context, model any, out any) error
	GetByID(ctx context.Context, model any, id uint, out any) error
	Create(ctx context.Context, model any) error
	UpdateByID(ctx context.Context, model any, id uint, updates map[string]any) error
	DeleteByID(ctx context.Context, model any, id uint) error
}

type ReservationRepository interface {
	MarkCheckIn(ctx context.Context, id uint) error
	MarkCheckOut(ctx context.Context, id uint) error
}

type ParkingTxRepository interface {
	MarkCheckOut(ctx context.Context, id uint) error
}

type Repositories struct {
	Auth        AuthRepository
	Crud        CrudRepository
	Reservation ReservationRepository
	ParkingTx   ParkingTxRepository
}

func New(db *gorm.DB) Repositories {
	return Repositories{
		Auth:        &authRepo{db: db},
		Crud:        &crudRepo{db: db},
		Reservation: &reservationRepo{db: db},
		ParkingTx:   &parkingTxRepo{db: db},
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

func (r *crudRepo) List(ctx context.Context, model any, out any) error {
	return r.db.WithContext(ctx).Model(model).Order("id DESC").Find(out).Error
}
func (r *crudRepo) GetByID(ctx context.Context, model any, id uint, out any) error {
	return r.db.WithContext(ctx).Model(model).First(out, id).Error
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

func (r *reservationRepo) MarkCheckIn(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&models.Reservation{}).Where("id = ?", id).Updates(map[string]any{"status": "checked_in", "actual_check_in": time.Now().UTC()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *reservationRepo) MarkCheckOut(ctx context.Context, id uint) error {
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

func (r *parkingTxRepo) MarkCheckOut(ctx context.Context, id uint) error {
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
