package env

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/cfg"
	"gitlab.com/g6834/team41/auth/internal/repositories"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Environment struct {
	C  *cfg.Config
	L  *logrus.Logger
	UR *repositories.UserRepository
	M  *Metrics
}

type Metrics struct {
	StartServiceCounter prometheus.Counter
	LoginCounter        prometheus.Counter
	LogoutCounter       prometheus.Counter
	ValidateCounter     prometheus.Counter
}

var e *Environment
var once sync.Once

const (
	HostAddress    = "HOST_ADDRESS"
	MetricsAddress = "METRICS_ADDRESS"
	ServiceName    = "auth-service"
	TracerName     = ServiceName
)

func E() *Environment {
	once.Do(func() {
		args := parseArgs()

		//init Environment
		e = &Environment{}

		e.L = logrus.New()

		//load config
		e.C = &cfg.Config{}
		if err := loadYaml(args["c"], e.C); err != nil {
			panic(err)
		}

		overrideСonfig := func(name string, value *string) {
			v := os.Getenv(name)
			if len(v) > 0 {
				*value = v
			}
		}
		overrideСonfig(HostAddress, &e.C.HostAddress)
		overrideСonfig(MetricsAddress, &e.C.MetricsAddress)

		//init Sentry
		err := sentry.Init(sentry.ClientOptions{
			Dsn:   e.C.SentryDSN,
			Debug: true,
		})
		if err != nil {
			panic(fmt.Errorf("sentry.Init: %w", err))
		}

		//init otel tracer (jaeger)
		tp, err := tracerProvider(ServiceName, e.C.JaegerCollector)
		if err != nil {
			panic(err)
		}
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.TraceContext{})

		//init prometheus
		e.M = &Metrics{}
		e.M.StartServiceCounter = promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "team41", Subsystem: "auth", Name: "startServiceCounter", Help: "Service start counter",
		})
		e.M.LoginCounter = promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "team41", Subsystem: "auth", Name: "loginCounter", Help: "Login endpoint request counter",
		})
		e.M.LogoutCounter = promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "team41", Subsystem: "auth", Name: "logoutCounter", Help: "Logout endpoint request counter",
		})
		e.M.ValidateCounter = promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "team41", Subsystem: "auth", Name: "validateCounter", Help: "Validate endpoint request counter",
		})

		// TODO: add repositories.UserRepository implementation.

		err = e.validate()
		if err != nil {
			sentry.CaptureException(err)
			panic(err)
		}

		//test
		sentry.CaptureMessage("Is it works?")
	})

	return e
}

// validate ensures that all required environment variables has been set.
func (E *Environment) validate() error {
	switch {
	case E.C.HostAddress == "":
		return errors.New("$" + HostAddress + " isn't set")
	case E.C.MetricsAddress == "":
		return errors.New("$" + MetricsAddress + " isn't set")
	case E.C.SentryDSN == "":
		return errors.New("sentry DSN isn't configured")
	}

	return nil
}

func (E *Environment) Tracer() trace.Tracer {
	return otel.Tracer(TracerName)
}

func tracerProvider(serviceName string, url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	return tp, nil
}
