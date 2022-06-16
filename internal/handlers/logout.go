package handlers

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type Logout struct{}

func (l Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Logout")
	defer span.End()

	env.E().M.LogoutCounter.Inc()

	// TODO:
	// logout
	// +redirect_uri

	sentry.CaptureException(errors.New("not yet implemented"))
	panic("not yet implemented")
}
