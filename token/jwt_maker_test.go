package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func generateRandomJWT(t *testing.T, username string, expTime time.Duration) string {
	signed, err := testJWTMaker.CreateToken(username, expTime)
	require.NoError(t, err)
	require.NotEqual(t, signed, "")
	return signed
}

func TestGenerateJWT(t *testing.T) {
	username := "sodiq38"
	signed := generateRandomJWT(t, username, time.Hour*2)
	t.Log(signed)
}

func TestValidateJWT(t *testing.T) {
	username := "olaniyi"
	signed := generateRandomJWT(t, username, time.Hour*2)

	payload, err := testJWTMaker.VerifyToken(signed)
	t.Log(payload)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotEmpty(t, payload.ID)
	require.Equal(t, payload.Username, username)
}

func TestExpiredToken(t *testing.T) {
	username := "sodiq"
	signed := generateRandomJWT(t, username, time.Second*2)

	time.Sleep(time.Second * 5)
	_, err := testJWTMaker.VerifyToken(signed)
	require.ErrorContains(t, err, jwt.ErrTokenExpired.Error())
}

func TestInvalidSigningMethod(t *testing.T) {
	HS384Signed := "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI5MjUwNjIsImlhdCI6MTc0MjkxNzg2MiwidXNlcm5hbWUiOiJzb2RpcTM4In0.TCXS4UzknXgGC0BlN_DrQ3xwvZcJEMGKYeDvbYht_xNVqpSjWc4-_Gx1yT4kCX88"
	t.Log("Checking for HS256 key")
	_, err := testJWTMaker.VerifyToken(HS384Signed)
	require.ErrorContains(t, err, jwt.ErrSignatureInvalid.Error())
}
