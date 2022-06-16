package handlers

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type Login struct{}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Login")
	defer span.End()

	env.E().M.LoginCounter.Inc()

	// TODO:
	// basic-авторизация
	// возвращает куки access и refresh, содержащие jwt-токены
	// +redirect_uri

	sentry.CaptureException(errors.New("not yet implemented"))
	panic("not yet implemented")
}
