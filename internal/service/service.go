package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"hotel/backend/internal/models"
	"hotel/backend/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Services struct {
	Auth        *AuthService
	Crud        *CrudService
	Reservation *ReservationService
	ParkingTx   *ParkingTxService
	Parking     *ParkingService
}

func New(repos repository.Repositories, sessionTTL time.Duration) Services {
	return Services{
		Auth:        &AuthService{AuthRepository: repos.Auth, sessionTTL: sessionTTL},
		Crud:        &CrudService{CrudRepository: repos.Crud},
		Reservation: &ReservationService{repos.Reservation},
		ParkingTx:   &ParkingTxService{repos.ParkingTx},
		Parking:     &ParkingService{repos.Parking},
	}
}

type AuthService struct {
	repository.AuthRepository
	sessionTTL time.Duration
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, time.Time, error) {
	user, err := s.GetUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, "", time.Time{}, errors.New("invalid credentials")
	}
	token, err := RandomToken()
	if err != nil {
		return nil, "", time.Time{}, err
	}
	expires := time.Now().UTC().Add(s.sessionTTL)
	if err := s.CreateSession(ctx, models.Session{ID: token, UserID: user.ID, ExpiresAt: expires}); err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expires, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	return s.DeleteSession(ctx, sessionID)
}

func (s *AuthService) Me(ctx context.Context, sessionID string) (*models.User, error) {
	return s.GetUserBySession(ctx, sessionID)
}

type CrudService struct{ repository.CrudRepository }

type ReservationService struct {
	repository.ReservationRepository
}

type ParkingTxService struct {
	repository.ParkingTxRepository
}

type ParkingService struct {
	repository.ParkingRepository
}

func RandomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf) + "-" + uuid.NewString(), nil
}

func IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
