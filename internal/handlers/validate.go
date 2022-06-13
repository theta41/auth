package handlers

import "net/http"

type Validate struct{}

func (v Validate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// по эндпоинту /i проверять jwt-куки на валидность и
	// возвращать данные залогиненного пользователя (пока только логин).
	// Предварительно jwt-токен должен быть провалидирован в мидлваре, разобран,
	// и структура пользователя должна быть передана в хэндлер через контекст

	panic("not yet implemented")
}
