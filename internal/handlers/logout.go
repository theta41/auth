package handlers

import "net/http"

type Logout struct{}

func (l Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// logout
	// +redirect_uri

	panic("not yet implemented")
}
