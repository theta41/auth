package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/g6834/team41/auth/internal/auth"
	"gitlab.com/g6834/team41/auth/internal/models"
	"gitlab.com/g6834/team41/auth/internal/repositories"
)

type JwtConfigInterface interface {
	GetJWTSecret() string
	GetJWTTTL() int
}

type Service struct {
	db repositories.UserRepository
	//TODO add logger?
	cfg JwtConfigInterface
	//TODO add tools? (auth.GetHash, auth.NewJWT)
}

func New(db repositories.UserRepository, cfg JwtConfigInterface) *Service {
	return &Service{
		db:  db,
		cfg: cfg,
	}
}

func (s *Service) Login(login, password string) (models.TokenPair, error) {
	user, err := s.db.GetUser(login)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("invalid login/password")
	}

	pass := auth.GetHash(password)
	if user.PasswordHash != pass {
		return models.TokenPair{}, fmt.Errorf("invalid login/password")
	}

	jwt, err := s.updateTokens(user)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("auth Login updateTokens error %w", err)
	}

	newTokens := models.TokenPair{
		AccessToken:  user.Token,
		RefreshToken: jwt,
	}
	return newTokens, nil
}

func (s *Service) Info(login string) (*models.User, error) {
	return s.db.GetUser(login)
}

func (s *Service) Validate(login string, tokens models.TokenPair) (models.TokenPair, error) {
	user, err := s.db.GetUser(login)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("auth Validate GetUser error: %w", err)
	}

	//TODO check JWT?
	if user.Token != tokens.AccessToken {
		return models.TokenPair{}, fmt.Errorf("invalid tokens pair")
	}

	jwt, err := s.updateTokens(user)
	if err != nil {
		return models.TokenPair{}, fmt.Errorf("auth Validate updateTokens error %w", err)
	}

	newTokens := models.TokenPair{
		AccessToken:  user.Token,
		RefreshToken: jwt,
	}
	return newTokens, nil
}

func (s *Service) updateTokens(user *models.User) (jwt string, err error) {
	user.Token = uuid.New().String()
	err = s.db.ChangeToken(user.Token, user.Login)
	if err != nil {
		return "", fmt.Errorf("change token error: %w", err)
	}

	expire := time.Now().Add(time.Duration(s.cfg.GetJWTTTL()) * time.Second)
	jwt, err = auth.NewJWT(user.Login, expire, s.cfg.GetJWTSecret())
	if err != nil {
		return "", fmt.Errorf("new JWT error: %w", err)
	}

	return jwt, nil
}
