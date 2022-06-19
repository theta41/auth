package main

import (
	"gitlab.com/g6834/team41/auth/internal/app"
	"gitlab.com/g6834/team41/auth/internal/env"
)

func main() {
	env.E().L.Info("starting service...")
	env.E().M.StartServiceCounter.Inc()

	defer env.OnStop()
	//TODO gracefull shutdown app

	a := app.NewApp()
	err := a.Run()
	if err != nil {
		env.E().L.Panic(err)
	}
}
