package handlers

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
	"net/http"
)

type Profiling struct{}

func (p Profiling) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	env.E().C.Profiling = !env.E().C.Profiling

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("{\"status\": %v}", env.E().C.Profiling)))
	if err != nil {
		sentry.CaptureException(err)
		env.E().L.Error(err)
	}
}
