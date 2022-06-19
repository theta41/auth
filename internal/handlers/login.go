package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"gitlab.com/g6834/team41/auth/internal/auth"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Login struct{}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Login")
	defer span.End()

	env.E().M.LoginCounter.Inc()

	w.Header().Add("Content-Type", "application/json")

	err := handle(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{}"))
		sentry.CaptureException(err)
		env.E().L.Error(err)
	}
}

func handle(w http.ResponseWriter, r *http.Request) error {
	// Parse body.
	req, err := parseRequest(r)
	if err != nil {
		return fmt.Errorf("error. parse request: %w", err)
	}

	// Get user from repository.
	user, err := env.E().UR.GetUser(req.Login)
	if err != nil {
		return fmt.Errorf("error. get user: %w", err)
	}

	pass := auth.GetHash(req.Password)

	if user.PasswordHash != pass {
		w.WriteHeader(http.StatusForbidden)
		_, err = w.Write([]byte("{}"))
		if err != nil {
			return fmt.Errorf("error. pass not equal: %w", err)
		}
		return nil
	}

	// Generate and save new refresh token.
	user.Token = uuid.New().String()
	err = env.E().UR.ChangeToken(user.Token, user.Login)
	if err != nil {
		return fmt.Errorf("Change token error: %w", err)
	}

	jwt, err := auth.NewJWT(user.Login, time.Now().Add(time.Duration(env.E().C.JWTTTL)*time.Second), env.E().C.JWTSecret)
	if err != nil {
		return fmt.Errorf("New JWT error error: %w", err)
	}

	resp := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		AccessToken:  user.Token,
		RefreshToken: jwt,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("Json marshal err: %w", err)
	}

	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("Write: %w", err)
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
