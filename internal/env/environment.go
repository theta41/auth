package env

import (
	"errors"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/cfg"
)

type Environment struct {
	C *cfg.Config
	L *logrus.Logger
}

var e *Environment
var once sync.Once

const (
	HostAddress = "HOST_ADDRESS"
)

func E() *Environment {
	once.Do(func() {
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
	})

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
