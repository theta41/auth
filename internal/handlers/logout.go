package handlers

import (
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/env"
)

type Logout struct{}

func (l Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Logout")
	defer span.End()

	env.E().M.LogoutCounter.Inc()

	w.Header().Set("Set-Cookie", "Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly")
}
