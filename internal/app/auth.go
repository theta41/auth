package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/handlers"
	"gitlab.com/g6834/team41/auth/internal/middlewares"
	"gitlab.com/g6834/team41/auth/internal/repositories"
	"net/http"
	"net/http/pprof"
)

type App struct {
	m  *chi.Mux
	ur repositories.UserRepository
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

	//start prometheus
	//TODO make it nicer later :)
	http.Handle(MetricsPath, promhttp.Handler())
	go http.ListenAndServe(env.E().C.MetricsAddress, nil)

	return http.ListenAndServe(env.E().C.HostAddress, a.m)
}

const (
	LoginPath    = "/login"
	LogoutPath   = "/logout"
	ValidatePath = "/validate"
	MetricsPath  = "/metrics"
	ProfilePath  = "/profile"
)

func (a *App) bindHandlers() {
	a.m.Handle(LoginPath, handlers.Login{})
	a.m.Handle(LogoutPath, handlers.Logout{})
	a.m.Handle(ValidatePath, handlers.Validate{})

	a.m.Handle(ProfilePath, handlers.Profiling{})

	a.m.Route("/debug/pprof", func(r chi.Router) {
		r.Use(middlewares.CheckProf)

		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
	})
}

func (a *App) registerMiddleware() {
	//a.m.Use(middleware.Logger)
	a.m.Use(middlewares.Logrus)
}
