package middlewares

import (
	"crypto/subtle"
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/auth"
	"gitlab.com/g6834/team41/auth/internal/repositories"
)

// BasicAuth copied from
// ref: https://www.alexedwards.net/blog/basic-authentication-in-go

func NewBasicAuth(db repositories.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			username, password, ok := r.BasicAuth()

			if ok {
				user, err := db.GetUser(username)
				if err == nil {

					// Prepare
					passwordHash := []byte(auth.GetHash(password))
					expectedPasswordHash := []byte(user.PasswordHash)

					// Constant time compare
					passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

					// Check
					if passwordMatch {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			// Unauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "{}", http.StatusForbidden)
		})
	}
}
