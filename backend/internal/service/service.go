package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"hotel/backend/internal/models"
	"hotel/backend/internal/repository"
)

type Services struct {
	Auth        *AuthService
	Crud        *CrudService
	Reservation *ReservationService
	ParkingTx   *ParkingTxService
}

func New(repos repository.Repositories, sessionTTL time.Duration) Services {
	return Services{
		Auth:        &AuthService{repo: repos.Auth, sessionTTL: sessionTTL},
		Crud:        &CrudService{repo: repos.Crud},
		Reservation: &ReservationService{repo: repos.Reservation},
		ParkingTx:   &ParkingTxService{repo: repos.ParkingTx},
	}
}

type AuthService struct {
	repo       repository.AuthRepository
	sessionTTL time.Duration
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, time.Time, error) {
	user, err := s.repo.GetUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, "", time.Time{}, errors.New("invalid credentials")
	}
	token, err := randomToken()
	if err != nil {
		return nil, "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(s.sessionTTL)
	if err := s.repo.CreateSession(ctx, models.Session{ID: token, UserID: user.ID, ExpiresAt: expires}); err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expires, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.repo.DeleteSession(ctx, sessionID)
}

func (s *AuthService) Me(ctx context.Context, sessionID string) (*models.User, error) {
	return s.repo.GetUserBySession(ctx, sessionID)
}

type CrudService struct{ repo repository.CrudRepository }

func (s *CrudService) List(ctx context.Context, model any, out any) error {
	return s.repo.List(ctx, model, out)
}
func (s *CrudService) GetByID(ctx context.Context, model any, id uint, out any) error {
	return s.repo.GetByID(ctx, model, id, out)
}
func (s *CrudService) Create(ctx context.Context, model any) error { return s.repo.Create(ctx, model) }
func (s *CrudService) UpdateByID(ctx context.Context, model any, id uint, updates map[string]any) error {
	return s.repo.UpdateByID(ctx, model, id, updates)
}
func (s *CrudService) DeleteByID(ctx context.Context, model any, id uint) error {
	return s.repo.DeleteByID(ctx, model, id)
}

type ReservationService struct {
	repo repository.ReservationRepository
}

func (s *ReservationService) CheckIn(ctx context.Context, id uint) error {
	return s.repo.MarkCheckIn(ctx, id)
}
func (s *ReservationService) CheckOut(ctx context.Context, id uint) error {
	return s.repo.MarkCheckOut(ctx, id)
}

type ParkingTxService struct {
	repo repository.ParkingTxRepository
}

func (s *ParkingTxService) CheckOut(ctx context.Context, id uint) error {
	return s.repo.MarkCheckOut(ctx, id)
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf) + "-" + uuid.NewString(), nil
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
