package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {

		// TODO:
		// basic-авторизация
		// возвращает куки access и refresh, содержащие jwt-токены
		// +redirect_uri

		w.Write([]byte("ok"))
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {

		// TODO:
		// logout
		// +redirect_uri

		w.Write([]byte("ok"))
	})

	r.Get("/i", func(w http.ResponseWriter, r *http.Request) {

		// TODO:
		// по эндпоинту /i проверять jwt-куки на валидность и
		// возвращать данные залогиненного пользователя (пока только логин).
		// Предварительно jwt-токен должен быть провалидирован в мидлваре, разобран,
		// и структура пользователя должна быть передана в хэндлер через контекст

		w.Write([]byte("ok"))
	})

	http.ListenAndServe(":3000", r)
}
