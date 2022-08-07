package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team41/auth/internal/middlewares"
)

type profilingConfig struct {
	enable bool
}

func (p profilingConfig) GetProfiling() bool {
	return p.enable
}

func TestCheckProf(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	cases := []struct {
		name              string
		profilingEnable   bool
		expectNextCounter int
		expectStatus      int
	}{
		{"profiling_disabled", false, 0, http.StatusForbidden},
		{"profiling_enabled", true, 1, http.StatusOK},
	}

	for _, tCase := range cases {
		tc := tCase
		t.Run(tc.name, func(tRun *testing.T) {

			var nextHandlerWasCalled int
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//tests may also goes here..
				w.WriteHeader(http.StatusOK)

				nextHandlerWasCalled++
			})

			cfg := profilingConfig{tc.profilingEnable}

			handlerCheckProf := middlewares.NewCheckProf(logger, cfg)(nextHandler)

			inRequest, _ := http.NewRequest("GET", "http://testing", nil)
			middlewareRecorder := httptest.NewRecorder()
			handlerCheckProf.ServeHTTP(middlewareRecorder, inRequest)

			assert.Equal(tRun, tc.expectNextCounter, nextHandlerWasCalled, "nextHandlerWasCalled")
			assert.Equal(tRun, tc.expectStatus, middlewareRecorder.Code)
		})
	}
}
