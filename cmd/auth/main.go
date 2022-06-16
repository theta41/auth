package main

import (
	"gitlab.com/g6834/team41/auth/internal/app"
	"gitlab.com/g6834/team41/auth/internal/env"
	"gitlab.com/g6834/team41/auth/internal/jsondb"
)

func main() {
	env.E().L.Info("starting service...")

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
