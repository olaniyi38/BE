package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olaniyi38/BE/token"
)

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordUpdatedAt time.Time `json:"password_updated_at"`
	FullName          string    `json:"full_name"`
}

func (server *Server) getUser(c *gin.Context) {
	payload, ok := c.MustGet(auth_payload_key).(*token.Payload)
	if !ok {
		respondWithErr(http.StatusUnauthorized, fmt.Errorf("bad request"), "no auth session found", c)
		return
	}

	user, err := server.store.GetUser(c, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithErr(http.StatusNotFound, fmt.Errorf("not found"), fmt.Sprintf("user of username %v not found", payload.Username), c)
			return
		}
		respondWithErr(http.StatusInternalServerError, err, "error getting user", c)
		return
	}

	response := userResponse{
		Username:          user.Username,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}
