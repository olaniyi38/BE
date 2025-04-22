package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	token "github.com/olaniyi38/BE/token"
	"github.com/olaniyi38/BE/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func createAuthCookie(token string, duration int) *http.Cookie {
	return &http.Cookie{
		Name:     auth_cookie_key,
		Value:    token,
		MaxAge:   duration,
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HttpOnly: true,
	}
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		addCookie     func(request *http.Request, maker token.Maker, duration int)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			addCookie: func(request *http.Request, maker token.Maker, duration int) {
				access_token, err := maker.CreateToken("username", time.Hour*2) //create token
				require.NoError(t, err)
				request.AddCookie(createAuthCookie(access_token, duration))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				response := gin.H{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, recorder.Code, http.StatusOK)
				require.Equal(t, response, gin.H{})
			},
		},
		{
			name:      "UNAUTHORIZED: COOKIE NOT FOUND",
			addCookie: func(request *http.Request, maker token.Maker, duration int) {},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				response := gin.H{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "UNAUTHORIZED: EXPIRED",
			addCookie: func(request *http.Request, maker token.Maker, duration int) {
				access_token, err := maker.CreateToken("username", -time.Hour*2) //create token
				require.NoError(t, err)
				request.AddCookie(createAuthCookie(access_token, duration))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				response := gin.H{}
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testSymmetricKey := util.RandomString(32)
	testMaker, err := token.NewPasetoMaker(testSymmetricKey)
	require.NoError(t, err)

	testAuthMiddleware := authMiddleware(testMaker)
	server := newTestServer(t, nil)

	server.router.GET("/auth", testAuthMiddleware, func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest("GET", "/auth", nil)
			cookieDuration := int(time.Now().Add(server.config.TokenDuration).Unix())
			testCase.addCookie(request, testMaker, cookieDuration) //add to request cookie
			server.router.ServeHTTP(recorder, request)             //run server
			testCase.checkResponse(t, recorder)                    //check response
		})
	}

}
