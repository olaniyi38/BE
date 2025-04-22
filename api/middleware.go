package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olaniyi38/BE/token"
)

const (
	auth_payload_key = "auth_payload"
)

func authMiddleware(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken, err := c.Cookie(auth_cookie_key)

		if err != nil {
			respondWithErr(http.StatusUnauthorized, err, "error authenticating user", c)
			c.Abort()
			return
		}

		payload, err := maker.VerifyToken(authToken)
		if err != nil {
			respondWithErr(http.StatusUnauthorized, err, "invalid token", c)
			c.Abort()
			return
		}

		c.Set(auth_payload_key, payload)
		c.Next()
	}
}
