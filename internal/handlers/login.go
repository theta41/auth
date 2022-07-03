package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gitlab.com/g6834/team41/auth/internal/handlers/util"
	"gitlab.com/g6834/team41/auth/internal/ports"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

func (l Login) handle(w http.ResponseWriter, r *http.Request) error {
	// Parse body.
	req, err := parseRequest(r)
	if err != nil {
		return fmt.Errorf("error. parse request: %w", err)
	}

	// Login within domain
	tokens, err := l.auth.Login(req.Login, req.Password)
	if err != nil {
		return fmt.Errorf("invalid login/password")
	}

	// Prepare response
	util.PutLoginToCookie(w, req.Login)
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

func parseRequest(r *http.Request) (*LoginRequest, error) {
	req := LoginRequest{}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &req)
	if err != nil {
		return nil, err
	}

	return &req, nil
}
