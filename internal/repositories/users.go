package repositories

import "gitlab.com/g6834/team41/auth/internal/models"

type UserRepository interface {
	// AddToken changing token of the user with given login.
	AddToken(token, login string) error

	// CheckToken validates, that the user with given login have specified token.
	CheckToken(token, login string) (bool, error)

	// GetUser selecting one user by login
	GetUser(login string) (models.User, error)
}
