package handlers

import (
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/handlers/util"
)

type Logout struct{}

// @Summary Logout
// @Description Logout
// @Accept json
// @Produce json
// @Success 200
// @Router /logout [get]
func (l Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Logout")
	defer span.End()

	env.E().M.LogoutCounter.Inc()

	//clear cookies
	util.ClearLoginCookie(w)
	util.ClearTokensCookie(w)

	w.Header().Set("Set-Cookie", "Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly")
}
