package middlewares

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/env"
	"net/http"
	"time"
)

func Logrus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			l := env.E().L.WithFields(logrus.Fields{
				"method":   r.Method,
				"path":     r.URL.Path,
				"status":   ww.Status(),
				"duration": time.Since(t1),
			})
			l.Info()
		}()

		next.ServeHTTP(w, r)
	})
}
