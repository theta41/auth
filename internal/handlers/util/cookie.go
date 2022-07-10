package util

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/auth/internal/models"
)

const (
	cookieLogin        = "Login"
	cookieAccessToken  = "AccessToken"
	cookieRefreshToken = "RefreshToken"
)

func mustCookie(r *http.Request, name string) string {
	v, err := r.Cookie(name)
	if err != nil || v == nil {
		logrus.Infof("missing cookie %s", name)
		return ""
	}

	// logrus.Debugf ?
	logrus.Infof("got cookie %v", v)
	return v.Value
}

func GetTokensFromCookie(r *http.Request) models.TokenPair {

	access := mustCookie(r, cookieAccessToken)
	refresh := mustCookie(r, cookieRefreshToken)

	return models.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func GetLoginFromCookie(r *http.Request) string {
	//return "test@example.org"

	login := mustCookie(r, cookieLogin)

	return login
}

func PutLoginToCookie(w http.ResponseWriter, loginValue string) {
	loginCookie := http.Cookie{
		Name:    cookieLogin,
		Value:   loginValue,
		Expires: time.Time{}.AddDate(9999, 0, 0), //learning cookies never expires
	}

	logrus.Print("put login to cookies", loginCookie)

	http.SetCookie(w, &loginCookie)
}

func ClearLoginCookie(w http.ResponseWriter) {
	loginCookie := http.Cookie{
		Name:     cookieLogin,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	logrus.Print("clear login cookie", loginCookie)

	http.SetCookie(w, &loginCookie)
}

func PutTokensToCookie(w http.ResponseWriter, tokens models.TokenPair) {
	access := http.Cookie{
		Name:    cookieAccessToken,
		Value:   tokens.AccessToken,
		Expires: time.Time{}.AddDate(9999, 0, 0), //learning cookies never expires
	}
	refresh := http.Cookie{
		Name:    cookieRefreshToken,
		Value:   tokens.RefreshToken,
		Expires: time.Time{}.AddDate(9999, 0, 0), //learning cookies never expires
	}

	logrus.Print("put tokens to cookies", access, refresh)

	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)
}

func ClearTokensCookie(w http.ResponseWriter) {
	access := http.Cookie{
		Name:     cookieAccessToken,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	refresh := http.Cookie{
		Name:     cookieRefreshToken,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	logrus.Print("clear tokens cookies", access, refresh)

	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)
}
