package handlers

import (
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/env"
)

type Validate struct{}

// @Summary Validate
// @Description Validate
// @Accept json
// @Produce json
// @Success 200
// @Failure 403
// @Failure 500
// @Router /validate [get]
func (v Validate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Validate")
	defer span.End()

	env.E().M.ValidateCounter.Inc()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[]"))
}
