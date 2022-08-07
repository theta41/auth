package ports

import (
	"gitlab.com/g6834/team41/auth/internal/models"
)

type Auth interface {
	Info(login string) (*models.User, error)
	CreateTokens(login string) (models.TokenPair, error)
	Validate(login string, tokens models.TokenPair) (models.TokenPair, error)
}
