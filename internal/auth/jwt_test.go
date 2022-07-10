package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOi02MjEzNTU5NjgwMCwibG9naW4iOiJ0ZXN0QHRlc3Qub3JnIn0.4C0307R0y9tDoXlq4NgG6ipA-MTddzXIUF7-42W9D2o"
	token, err := NewJWT("test@test.org", time.Time{}, "the-secret")
	assert.NoError(t, err)
	assert.Equal(t, expected, token)
}
