package handlers

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type Profiling struct{}

// @Summary Profile
// @Description Profile
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /profile [get]
func (p Profiling) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	env.E().C.Profiling = !env.E().C.Profiling

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("{\"status\": %v}", env.E().C.Profiling)))
	if err != nil {
		sentry.CaptureException(err)
		env.E().L.Error(err)
	}
}
