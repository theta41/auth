package handlers

import "net/http"

type Login struct{}

func (l Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// basic-авторизация
	// возвращает куки access и refresh, содержащие jwt-токены
	// +redirect_uri

	panic("not yet implemented")
}
