package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/handlers"
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
)

func (a *App) bindHandlers() {
	a.m.Handle(LoginPath, handlers.Login{})
	a.m.Handle(LogoutPath, handlers.Logout{})
	a.m.Handle(ValidatePath, handlers.Validate{})

	a.m.Route("/debug/pprof", func(r chi.Router) {
		r.Use(CheckProf)

		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
	})
}

func CheckProf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !env.E().C.Profiling {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("{}"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (a *App) registerMiddleware() {
	a.m.Use(middleware.Logger)
}
