package api

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/token"
	"github.com/olaniyi38/BE/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		PasetoSymmetricKey: util.RandomString(32),
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func attachAuthToRequest(t *testing.T, request *http.Request, tokenMaker token.Maker, userName string) {
	token, err := tokenMaker.CreateToken(userName, time.Minute*15)
	require.NoError(t, err)
	request.AddCookie(createAuthCookie(token, 10000000000))
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
