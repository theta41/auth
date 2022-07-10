package handlers

import (
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/handlers/util"
)

type Info struct{}

// @Summary Info
// @Description Info
// @Accept json
// @Produce json
// @Success 200
// @Failure 403
// @Failure 500
// @Router /i [get]
func (i Info) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Validate")
	defer span.End()

	env.E().M.ValidateCounter.Inc()

	login := util.GetLoginFromCookie(r)
	env.E().L.Infof("login %v", login) //debug

	tokens := util.GetTokensFromCookie(r)
	env.E().L.Infof("tokens %v", tokens) //debug

	if tokens.AccessToken == "" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("{}"))
		return
	}

	// TODO find user by AccessToken ?

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("[]"))
}
