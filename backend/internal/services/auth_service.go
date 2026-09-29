package services

import (
	"errors"
	"log/slog"
	"time"
	"uuid"

	"github.com/golang-jwt/jwt/v5"
	"github.com/williamf6894/VB-Events/internal/config"
	"github.com/williamf6894/VB-Events/internal/models"
	"github.com/williamf6894/VB-Events/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService struct {
	repo *repository.ParticipantRepository
	cfg  config.Config
}

func NewAuthService(repo *repository.ParticipantRepository, cfg config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func validParticipant(p *models.Participant) bool {
	return p.Name != "" && p.Email != "" && len(p.Password) >= minPasswordLength
}

func (s *AuthService) Create(participant *models.Participant) error {
	if !validParticipant(participant) {
		return ErrInvalidParticipant
	}

	existing, err := s.repo.FindByEmail(participant.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return ErrEmailTaken
	}

	hash, err := hashPassword(participant.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return err
	}
	participant.Password = hash

	return s.repo.Create(participant)
}

func (s *AuthService) Login(req LoginRequest) (*TokenResponse, error) {
	participant, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(participant.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(participant)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{Token: token}, nil
}

func (s *AuthService) GenerateToken(participant *models.Participant) (string, error) {
	if s.cfg.JWTSecret == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}

	claims := jwt.MapClaims{
		"user_id": participant.ID.String(),
		"email":   participant.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil(), errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil(), errors.New("invalid claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil(), errors.New("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil(), errors.New("invalid user_id format")
	}

	return userID, nil
}
