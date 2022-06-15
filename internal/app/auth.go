package app

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/handlers"

	"net/http"
)

type App struct {
	m *chi.Mux
}

func NewApp() *App {
	a := &App{
		m: chi.NewRouter(),
	}

	return a
}

func (a *App) Run() error {
	a.registerMiddleware()
	a.bindHandlers()

	addr := fmt.Sprintf("%s:%s", env.E().C.HostAddress, env.E().C.HostPort)
	return http.ListenAndServe(addr, a.m)
}

const (
	LoginPath    = "/login"
	LogoutPath   = "/logout"
	ValidatePath = "/validate"
)

func (a *App) bindHandlers() {
	a.m.Handle(LoginPath, handlers.Login{})
	a.m.Handle(LogoutPath, handlers.Logout{})
	a.m.Handle(ValidatePath, handlers.Validate{})
}

func (a *App) registerMiddleware() {
	a.m.Use(middleware.Logger)
}
