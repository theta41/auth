package app

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/g6834/team41/auth/internal/domain/auth"
	"gitlab.com/g6834/team41/auth/internal/env"
	authgrpc "gitlab.com/g6834/team41/auth/internal/grpc"
	"gitlab.com/g6834/team41/auth/internal/handlers"
	"gitlab.com/g6834/team41/auth/internal/middlewares"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/g6834/team41/auth/docs"
)

type App struct {
	m  *chi.Mux
	ds *auth.Service
}

func NewApp() *App {
	a := &App{
		m:  chi.NewRouter(),
		ds: auth.New(env.E().UR, env.E().C),
	}

	return a
}

func (a *App) Run() error {
	a.registerMiddleware()
	a.bindHandlers()

	//start prometheus
	http.Handle(MetricsPath, promhttp.Handler())
	go http.ListenAndServe(env.E().C.MetricsAddress, nil)

	go authgrpc.StartServer(env.E().C.GrpcAddress, a.ds)

	env.E().L.Infof("auth RAPI ListenAndServe at %v", env.E().C.HostAddress)
	return http.ListenAndServe(env.E().C.HostAddress, a.m)
}

const (
	LoginPath    = "/login"
	LogoutPath   = "/logout"
	ValidatePath = "/validate"
	InfoPath     = "/i"
	MetricsPath  = "/metrics"
	ProfilePath  = "/profile"
)

func (a *App) bindHandlers() {
	a.m.Handle(
		LoginPath,
		middlewares.NewBasicAuth(env.E().UR)(handlers.NewLogin(a.ds)),
	)
	a.m.Handle(
		LogoutPath,
		middlewares.NewBasicAuth(env.E().UR)(handlers.Logout{}),
	)

	a.m.Handle(ValidatePath, handlers.Validate{})
	a.m.Handle(InfoPath, handlers.Info{})

	a.m.Handle(ProfilePath, handlers.Profiling{})

	a.m.Route("/debug/pprof", func(r chi.Router) {
		r.Use(middlewares.NewCheckProf(env.E().L, env.E().C))

		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
	})

	bindSwagger(a.m)
}

func (a *App) registerMiddleware() {
	//a.m.Use(middleware.Logger)
	a.m.Use(middlewares.NewLogrus(env.E().L))
	//a.m.Use(middlewares.NewBasicAuth(env.E().UR))
}

func bindSwagger(r *chi.Mux) {
	r.Route("/swagger", func(r chi.Router) {
		r.HandleFunc("/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost"+env.E().C.HostAddress+"/swagger/doc.json"),
		))
	})
}
