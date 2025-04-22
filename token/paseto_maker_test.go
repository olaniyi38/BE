package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func generateToken(t *testing.T, username string, duration time.Duration) string {
	token, err := testPasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}

func TestGenerateToken(t *testing.T) {
	generateToken(t, "sodiq", time.Hour*2)
}

func TestVerifyToken(t *testing.T) {
	username := "sodiq 38"

	token := generateToken(t, username, time.Hour*2)

	payload, err := testPasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)
	require.Equal(t, payload.Username, username)
}

func TestExpiredPasetoToken(t *testing.T) {
	username := "sodiq 38"

	token := generateToken(t, username, time.Second*2)
	time.Sleep(time.Second * 3)

	_, err := testPasetoMaker.VerifyToken(token)
	require.ErrorContains(t, err, ErrTokenExpired.Error())

}
