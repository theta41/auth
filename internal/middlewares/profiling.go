package middlewares

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type ProfilingConfigInterface interface {
	GetProfiling() bool
}

func NewCheckProf(
	logger *logrus.Logger,
	cfg ProfilingConfigInterface,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.GetProfiling() {
				w.WriteHeader(http.StatusForbidden)
				_, err := w.Write([]byte("{}"))
				if err != nil {
					sentry.CaptureException(err)
					logger.Error(err)
				}
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
