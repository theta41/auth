package env

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	StartServiceCounter prometheus.Counter
	LoginCounter        prometheus.Counter
	LogoutCounter       prometheus.Counter
	ValidateCounter     prometheus.Counter
}

const (
	team41 = "team41"
)

func newMetrics(name string) *Metrics {
	return &Metrics{
		StartServiceCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: team41, Subsystem: name, Name: "startServiceCounter", Help: "Service start counter",
		}),

		LoginCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: team41, Subsystem: name, Name: "loginCounter", Help: "Login endpoint request counter",
		}),

		LogoutCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: team41, Subsystem: name, Name: "logoutCounter", Help: "Logout endpoint request counter",
		}),

		ValidateCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: team41, Subsystem: name, Name: "validateCounter", Help: "Validate endpoint request counter",
		}),
	}
}
