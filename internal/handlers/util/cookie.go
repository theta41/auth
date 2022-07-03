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
	loginCookie := http.Cookie{Name: cookieLogin}

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
	access := http.Cookie{Name: cookieAccessToken}
	refresh := http.Cookie{Name: cookieRefreshToken}

	logrus.Print("clear tokens cookies", access, refresh)

	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)
}
