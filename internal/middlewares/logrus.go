package middlewares

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewLogrus(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				l := logger.WithFields(logrus.Fields{
					"method":        r.Method,
					"path":          r.URL.Path,
					"status":        ww.Status(),
					"bytes_written": ww.BytesWritten(),
					"header":        ww.Header(),
					"duration":      time.Since(t1),
				})
				l.Info()
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
