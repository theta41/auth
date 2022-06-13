package main

import (
	"gitlab.com/g6834/team41/auth/internal/app"
	"gitlab.com/g6834/team41/auth/internal/env"
)

func main() {
	env.E().L.Info("starting service...")

	a := app.NewApp()
	err := a.Run()
	if err != nil {
		env.E().L.Panic(err)
	}
}
