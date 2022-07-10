package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gitlab.com/g6834/team41/auth/internal/middlewares"
)

func TestLoggerMiddleware(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetReportCaller(true)

	var nextHandlerWasCalled int
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//tests may also goes here..

		nextHandlerWasCalled++
	})

	handlerLogger := middlewares.NewLogrus(logger)(nextHandler)

	inRequest, _ := http.NewRequest("GET", "http://testing", nil)
	middlewareRecorder := httptest.NewRecorder()
	handlerLogger.ServeHTTP(middlewareRecorder, inRequest)

	require.Greater(t, nextHandlerWasCalled, 0, "next handler wasn't called")
}
