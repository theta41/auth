package handlers

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"gitlab.com/g6834/team41/auth/internal/env"
)

type Validate struct{}

func (v Validate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, span := env.E().Tracer().Start(r.Context(), "Validate")
	defer span.End()

	env.E().M.ValidateCounter.Inc()

	// TODO:
	// по эндпоинту /i проверять jwt-куки на валидность и
	// возвращать данные залогиненного пользователя (пока только логин).
	// Предварительно jwt-токен должен быть провалидирован в мидлваре, разобран,
	// и структура пользователя должна быть передана в хэндлер через контекст

	sentry.CaptureException(errors.New("not yet implemented"))
	panic("not yet implemented")
}
