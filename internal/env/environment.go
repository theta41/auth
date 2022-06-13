package env

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/cfg"
	"gitlab.com/g6834/team41/auth/internal/repositories"
	"os"
)

type Environment struct {
	C  *cfg.Config
	L  *logrus.Logger
	UR *repositories.UserRepository
}

var e *Environment

const (
	HostAddress = "HOST_ADDRESS"
)

func E() *Environment {
	if e == nil {
		e = &Environment{}

		e.L = logrus.New()
		e.C = &cfg.Config{
			HostAddress: os.Getenv(HostAddress),
		}

		// TODO: add repositories.UserRepository implementation.

		err := e.validate()
		if err != nil {
			panic(err)
		}
	}

	return e
}

// validate ensures that all required environment variables has been set.
func (E *Environment) validate() error {
	switch {
	case E.C.HostAddress == "":
		return errors.New("$" + HostAddress + " isn't set")
	}

	return nil
}
