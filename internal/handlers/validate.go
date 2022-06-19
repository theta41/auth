package handlers

import (
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/env"
)

type Validate struct{}

func (v Validate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Validate")
	defer span.End()

	env.E().M.ValidateCounter.Inc()

	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("[]"))
}
