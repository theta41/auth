package mongo

import (
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/auth"
	"gitlab.com/g6834/team41/auth/internal/models"
	"gitlab.com/g6834/team41/auth/internal/repositories"
)

const (
	hardcoded_DEFAULT_TEST_USER          = "test123"
	hardcoded_DEFAULT_TEST_USER_PASSWORD = "qwerty"
)

func BuildDefaultsIfNeeded(db repositories.UserRepository, logger *logrus.Logger) {

	if _, err := db.GetUser(hardcoded_DEFAULT_TEST_USER); err == nil {
		return // nothing to do
	}

	err := db.AddUser(models.User{
		Login:        hardcoded_DEFAULT_TEST_USER,
		PasswordHash: auth.GetHash(hardcoded_DEFAULT_TEST_USER_PASSWORD),
		RegisterDate: time.Now(),
	})

	if err != nil {
		logger.Errorf("can't create DEFAULT_TEST_USER: %s", err.Error())
	}
}
