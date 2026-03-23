package service

import (
	"errors"
	"time"

	"github.com/MirMonajir244/BookMySlot/internal/config"
	"github.com/MirMonajir244/BookMySlot/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type CoachRepo interface {
	Create(coach *models.Coach) error
	FindByEmail(email string) (*models.Coach, error)
	FindByID(id uint) (*models.Coach, error)
}

type AuthService struct {
	userRepo  UserRepo
	coachRepo CoachRepo
	cfg       *config.Config
}

func NewAuthService(userRepo UserRepo, coachRepo CoachRepo, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		coachRepo: coachRepo,
		cfg:       cfg,
	}
}

func (s *AuthService) RegisterUser(name, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("email already registered")
	}

	return user, nil
}

func (s *AuthService) RegisterCoach(name, email, password string) (*models.Coach, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	coach := &models.Coach{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := s.coachRepo.Create(coach); err != nil {
		return nil, errors.New("email already registered")
	}

	return coach, nil
}

func (s *AuthService) LoginUser(email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(user.ID, "user")
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) LoginCoach(email, password string) (*models.Coach, string, error) {
	coach, err := s.coachRepo.FindByEmail(email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(coach.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(coach.ID, "coach")
	if err != nil {
		return nil, "", err
	}

	return coach, token, nil
}

func (s *AuthService) generateToken(id uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(time.Duration(s.cfg.JWTExpiry) * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
