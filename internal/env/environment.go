package env

import (
	"errors"
	"gitlab.com/g6834/team41/auth/internal/mongo"
	"gitlab.com/g6834/team41/auth/internal/repositories"
	"os"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/cfg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Environment struct {
	C  *cfg.Config
	L  *logrus.Logger
	M  *Metrics
	UR repositories.UserRepository
}

var e *Environment
var once sync.Once

const (
	HostAddress    = "HOST_ADDRESS"
	MetricsAddress = "METRICS_ADDRESS"
	JWTSecret      = "JWT_SECRET"
	DBPassword     = "DB_PASSWORD"
	DBLogin        = "DB_LOGIN"
	ServiceName    = "auth"
	TracerName     = ServiceName
)

func E() *Environment {
	once.Do(func() {
		// Prepare.
		args := parseArgs()
		configYamlFilename := args["c"]

		// Init Environment.
		e = &Environment{}

		// Configure logger
		e.L = logrus.New()
		e.L.SetFormatter(&logrus.JSONFormatter{})
		e.L.SetReportCaller(true)

		e.loadConfig(configYamlFilename)

		initSentry(e.C.SentryDSN)
		initJaeger("team41-"+ServiceName, e.C.JaegerCollector)

		// Init prometheus counters.
		e.M = newMetrics(ServiceName)

		// Validate environment state.
		err := e.validate()
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		e.L.Info("Connecting to database...")
		e.UR, err = mongo.NewUsers(e.C.DB.Login, e.C.DB.Password, e.C.DB.Address, e.C.DB.Name, e.C.DB.Port)
		if err != nil {
			e.L.Panic(err)
		}
	})

	return e
}

func OnStop() {
	sentry.Flush(1 * time.Second)
}

// validate ensures that all required environment variables has been set.
func (E *Environment) validate() error {
	toCheck := []struct {
		Name  string
		Value string
	}{
		{HostAddress, E.C.HostAddress},
		{MetricsAddress, E.C.MetricsAddress},
		{JWTSecret, E.C.JWTSecret},
		{DBPassword, E.C.DB.Password},
		{DBLogin, E.C.DB.Login},
		{"sentry DSN", E.C.SentryDSN},
		{"jaeger collector", E.C.JaegerCollector},
	}

	return checkEmpty(toCheck...)
}

func checkEmpty(s ...struct {
	Name  string
	Value string
}) error {
	for i := range s {
		if s[i].Value == "" {
			return errors.New("$" + s[i].Name + " isn't set")
		}
	}

	return nil
}

func (E *Environment) Tracer() trace.Tracer {
	return otel.Tracer(TracerName)
}

func (E *Environment) loadConfig(filename string) {
	var err error
	E.C, err = cfg.NewConfig(filename)
	if err != nil {
		panic(err)
	}

	getEnv := func(name, old string) string {
		if value, ok := os.LookupEnv(name); ok {
			return value
		}
		return old
	}

	E.C.HostAddress = getEnv(HostAddress, E.C.HostAddress)
	E.C.MetricsAddress = getEnv(MetricsAddress, E.C.MetricsAddress)
	E.C.JWTSecret = getEnv(JWTSecret, E.C.JWTSecret)
	E.C.DB.Password = getEnv(DBPassword, E.C.DB.Password)
	E.C.DB.Login = getEnv(DBLogin, E.C.DB.Login)
}
