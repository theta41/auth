package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/g6834/team41/auth/internal/middlewares"
	"gitlab.com/g6834/team41/auth/internal/models"
)

type stubdb struct{}

func (stub stubdb) ChangeToken(token, login string) error {
	return nil
}
func (stub stubdb) GetUser(login string) (*models.User, error) {
	return &models.User{
		PasswordHash: "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5",
	}, nil
}

func TestBasicAuth(t *testing.T) {
	cases := []struct {
		name              string
		header            http.Header
		expectNextCounter int
		expectStatus      int
	}{
		{"empty_header", http.Header{}, 0, http.StatusForbidden},
		{"valid_basic_auth", http.Header{"Authorization": {"Basic dGVzdDEyMzpxd2VydHk="}}, 1, http.StatusOK},
	}

	for _, tCase := range cases {
		tc := tCase
		t.Run(tc.name, func(t *testing.T) {

			var nextHandlerWasCalled int
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				//tests may also goes here..
				w.WriteHeader(http.StatusOK)

				nextHandlerWasCalled++
			})

			handlerBasicAuth := middlewares.NewBasicAuth(&stubdb{})(nextHandler)

			request := &http.Request{
				Method: "GET",
				URL:    &url.URL{Path: "http:/testing.basic.auth/"},
				Header: tc.header,
			}
			middlewareRecorder := httptest.NewRecorder()
			handlerBasicAuth.ServeHTTP(middlewareRecorder, request)

			assert.Equal(t, tc.expectNextCounter, nextHandlerWasCalled, "nextHandlerWasCalled")
			assert.Equal(t, tc.expectStatus, middlewareRecorder.Code)
		})
	}
}
