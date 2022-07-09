package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/g6834/team41/auth/internal/models"
)

func TestPutLogiCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	PutLoginToCookie(recorder, "a@a.org")

	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}

	login, err := request.Cookie(cookieLogin)
	require.NoError(t, err)
	require.NotNil(t, login)
	require.Equal(t, "a@a.org", login.Value)
}

func TestClearLoginCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	ClearLoginCookie(recorder)

	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}

	login, err := request.Cookie(cookieLogin)
	require.NoError(t, err)
	require.NotNil(t, login)
	require.Empty(t, login.Value)
}

func TestPutTokensCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	PutTokensToCookie(recorder, models.TokenPair{
		AccessToken:  "123",
		RefreshToken: "456",
	})

	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}

	cookieAccess, err := request.Cookie(cookieAccessToken)
	require.NoError(t, err)
	require.NotNil(t, cookieAccess)
	require.Equal(t, "123", cookieAccess.Value)

	cookieRefresh, err := request.Cookie(cookieRefreshToken)
	require.NoError(t, err)
	require.NotNil(t, cookieRefresh)
	require.Equal(t, "456", cookieRefresh.Value)
}

func TestClearTokensCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	ClearTokensCookie(recorder)

	request := &http.Request{Header: http.Header{"Cookie": recorder.Header()["Set-Cookie"]}}

	cookieAccess, err := request.Cookie(cookieAccessToken)
	require.NoError(t, err)
	require.NotNil(t, cookieAccess)
	require.Empty(t, cookieAccess.Value)

	cookieRefresh, err := request.Cookie(cookieRefreshToken)
	require.NoError(t, err)
	require.NotNil(t, cookieRefresh)
	require.Empty(t, cookieRefresh.Value)
}
