package env

import (
	"errors"
	"log"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/cfg"
	"gitlab.com/g6834/team41/auth/internal/repositories"

	"github.com/joho/godotenv"
)

type Environment struct {
	C  *cfg.Config
	L  *logrus.Logger
	UR *repositories.UserRepository
}

var e *Environment
var once sync.Once

const (
	HostAddress = "HOST_ADDRESS"
	HostPort    = "HOST_PORT"
)

func E() *Environment {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	once.Do(func() {
		e = &Environment{}

		e.L = logrus.New()
		e.C = &cfg.Config{
			HostAddress: os.Getenv(HostAddress),
			HostPort:    os.Getenv(HostPort),
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
