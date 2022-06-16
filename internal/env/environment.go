package env

import (
	"errors"
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
	C *cfg.Config
	L *logrus.Logger
	M *Metrics
}

var e *Environment
var once sync.Once

const (
	HostAddress    = "HOST_ADDRESS"
	MetricsAddress = "METRICS_ADDRESS"
	ServiceName    = "auth"
	TracerName     = ServiceName
)

func E() *Environment {
	once.Do(func() {
		//prepare
		args := parseArgs()
		configYamlFilename := args["c"]

		//init Environment
		e = &Environment{}
		e.L = logrus.New()
		e.loadConfig(configYamlFilename)

		initSentry(e.C.SentryDSN)
		initJaeger(ServiceName, e.C.JaegerCollector)

		//init prometheus counters
		e.M = newMetrics(ServiceName)

		//validate environment state
		err := e.validate()
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		//test
		sentry.CaptureMessage("Is it works?")
	})

	return e
}

func OnStop() {
	sentry.Flush(1 * time.Second)
}

// validate ensures that all required environment variables has been set.
func (E *Environment) validate() error {
	switch {
	case E.C.HostAddress == "":
		return errors.New("$" + HostAddress + " isn't set")
	case E.C.MetricsAddress == "":
		return errors.New("$" + MetricsAddress + " isn't set")
	case E.C.SentryDSN == "":
		return errors.New("sentry DSN isn't set")
	case E.C.JaegerCollector == "":
		return errors.New("jaeger collector isn't set")
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
}
