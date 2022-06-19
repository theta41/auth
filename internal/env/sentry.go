package env

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

func initSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: true,
	})
	if err != nil {
		panic(fmt.Errorf("sentry.Init: %w", err))
	}
}
