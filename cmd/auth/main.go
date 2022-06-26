package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/app"
	"gitlab.com/g6834/team41/auth/internal/env"
)

func main() {
	env.E().L.SetFormatter(&logrus.JSONFormatter{})
	env.E().L.SetReportCaller(true)

	env.E().L.Info("starting service...")

	defer env.OnStop()
	//TODO gracefull shutdown app

	a := app.NewApp()
	err := a.Run()
	if err != nil {
		env.E().L.Panic(err)
	}
}
