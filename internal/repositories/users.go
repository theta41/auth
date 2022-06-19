package repositories

import "gitlab.com/g6834/team41/auth/internal/models"

type UserRepository interface {
	// ChangeToken changing token of the user with given login.
	ChangeToken(token, login string) error

	// GetUser selecting one user by login
	GetUser(login string) (*models.User, error)
}
