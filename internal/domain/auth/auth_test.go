package auth_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
	"gitlab.com/g6834/team41/auth/internal/domain/auth"
	"gitlab.com/g6834/team41/auth/internal/models"
)

type AuthTestSuite struct {
	suite.Suite

	auth *auth.Service
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (a *AuthTestSuite) SetupSuite() {
	a.auth = auth.New(a, a)
}

// ConfigInterface ////
func (a *AuthTestSuite) GetJWTSecret() string {
	return "the-secret"
}
func (a *AuthTestSuite) GetJWTTTL() int {
	return 123
}

// UserRepository ////
func (a *AuthTestSuite) ChangeToken(token, login string) error {
	if login == "broken-change-token" {
		return errors.New(login)
	}
	return nil
}
func (a *AuthTestSuite) GetUser(login string) (*models.User, error) {
	if login == "" {
		return nil, errors.New("no login")
	}
	return &models.User{
		Login:        login,
		Token:        "Token",
		PasswordHash: "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5", //qwerty
	}, nil
}
func (a *AuthTestSuite) AddUser(user models.User) error { return nil }
func (a *AuthTestSuite) DeleteUser(login string) error  { return nil }

// Tests /////

func (a *AuthTestSuite) TestCreateTokens() {

	cases := []struct {
		name                  string
		login                 string
		expectAccessTokenLen  int
		expectRefreshTokenLen int
		expectError           bool
	}{
		{"no_login", "", 0, 0, true},
		{"error_update_token", "broken-change-token", 0, 0, true},
		{"good_login", "test", 36, 125, false},
	}

	for _, tCase := range cases {
		tc := tCase
		a.Run(tc.name, func() {

			tokens, err := a.auth.CreateTokens(tc.login)

			if tc.expectError {
				assert.Error(a.T(), err)
			} else {
				assert.NoError(a.T(), err)
			}

			assert.Len(a.T(), tokens.AccessToken, tc.expectAccessTokenLen, "AccessToken")
			assert.Len(a.T(), tokens.RefreshToken, tc.expectRefreshTokenLen, "RefreshToken")
		})
	}
}

func (a *AuthTestSuite) TestInfo() {
	cases := []struct {
		name        string
		login       string
		expectUser  bool
		expectError bool
	}{
		{"no_login", "", false, true},
		{"good_login", "test", true, false},
	}

	for _, tCase := range cases {
		tc := tCase
		a.Run(tc.name, func() {

			user, err := a.auth.Info(tc.login)

			if tc.expectError {
				assert.Error(a.T(), err)
			} else {
				assert.NoError(a.T(), err)
			}

			if tc.expectUser {
				assert.NotNil(a.T(), user)
				assert.NotEmpty(a.T(), *user)
			} else {
				assert.Nil(a.T(), user)
			}
		})
	}
}

func (a *AuthTestSuite) TestValidate() {
	cases := []struct {
		name        string
		login       string
		tokens      models.TokenPair
		expextError bool
	}{
		{"no_login", "", models.TokenPair{}, true},
		{"invalid_token", "test", models.TokenPair{AccessToken: "123"}, true},
		{"error_update_token", "broken-change-token", models.TokenPair{AccessToken: "Token"}, true},
		{"good_login", "test", models.TokenPair{AccessToken: "Token"}, false},
	}

	for _, tCase := range cases {
		tc := tCase
		a.Run(tc.name, func() {

			tokens, err := a.auth.Validate(tc.login, tc.tokens)

			if tc.expextError {
				require.Error(a.T(), err)
				assert.Empty(a.T(), tokens)
			} else {
				require.NoError(a.T(), err)
				assert.NotEmpty(a.T(), tokens.AccessToken)
				assert.NotEmpty(a.T(), tokens.RefreshToken)
				assert.NotEqual(a.T(), tc.tokens.AccessToken, tokens.AccessToken)
				assert.NotEqual(a.T(), tc.tokens.RefreshToken, tokens.RefreshToken)
			}
		})
	}
}
