package middlewares

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
	"net/http"
)

func CheckProf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !env.E().C.Profiling {
			w.WriteHeader(http.StatusForbidden)
			_, err := w.Write([]byte("{}"))
			if err != nil {
				sentry.CaptureException(err)
				env.E().L.Error(err)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
