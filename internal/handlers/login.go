package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/handlers/util"
	"gitlab.com/g6834/team41/auth/internal/models"
	"gitlab.com/g6834/team41/auth/internal/ports"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type LoginRequest struct {
	Login    string `json:"login" excample:"t@t.org"`
	Password string `json:"password" excample:"123456"`
}

type Login struct {
	auth ports.Auth
}

func NewLogin(auth ports.Auth) Login {
	return Login{
		auth: auth,
	}
}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Login")
	defer span.End()

	env.E().M.LoginCounter.Inc()

	w.Header().Add("Content-Type", "application/json")

	err := l.handle(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{}"))
		sentry.CaptureException(err)
		env.E().L.Error(err)
	}
}

// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Param task body LoginRequest true "Login"
// @Success 200
// @Failure 500
// @Router /login [get]
func (l Login) handle(w http.ResponseWriter, r *http.Request) error {

	username := r.Context().Value(models.CtxUsername{}).(string)

	tokens, err := l.auth.CreateTokens(username)
	if err != nil {
		return fmt.Errorf("invalid login/password")
	}

	util.PutLoginToCookie(w, username)
	util.PutTokensToCookie(w, tokens)

	resp := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("http.ResponseWriter write error: %w", err)
	}

	return nil
}
