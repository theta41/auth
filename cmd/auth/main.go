package main

import (
	"time"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/app"
	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/jsondb"
)

func main() {
	env.E().L.Info("starting service...")
	env.E().M.StartServiceCounter.Inc()

	defer sentry.Flush(1 * time.Second)
	//TODO gracefull shutdown app

	ur, err := jsondb.NewJsonUsers("example.json")
	if err != nil {
		env.E().L.Panic(err)
	}

	a := app.NewApp(ur)
	err = a.Run()
	if err != nil {
		env.E().L.Panic(err)
	}
}
