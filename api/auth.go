package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/util"
)

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		PasswordUpdatedAt: user.PasswordUpdatedAt,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt,
	}

}

type loginRequest struct {
	Username string `json:"username" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=1"`
}

func (server *Server) Login(c *gin.Context) {
	var request loginRequest
	if err := c.BindJSON(&request); err != nil {
		respondWithErr(http.StatusBadRequest, err, "invalid request", c)
		return
	}

	user, err := server.store.GetUser(c, request.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithErr(http.StatusNotFound, fmt.Errorf("not found"), fmt.Sprintf("user of username %v not found", request.Username), c)
			return
		}
		respondWithErr(http.StatusInternalServerError, err, "error getting user", c)
		return
	}

	passwordValid := util.CheckPassword(request.Password, user.Password)
	if passwordValid != nil {
		respondWithErr(http.StatusUnauthorized, fmt.Errorf("invalid password or username"), "", c)
		return
	}

	token, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenDuration)
	if err != nil {
		respondWithErr(http.StatusInternalServerError, fmt.Errorf("internal"), "error login user", c)
		return
	}

	cookieDuration := time.Now().Add(server.config.TokenDuration)
	c.SetCookie(auth_cookie_key, token, int(cookieDuration.Unix()), "/", "localhost", true, true)

	response := newUserResponse(user)
	c.JSON(http.StatusOK, response)

}

type signUpRequest struct {
	Username string `json:"username" binding:"required,min=2"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=2"`
}

func (server *Server) SignUp(c *gin.Context) {
	var request signUpRequest
	if err := c.BindJSON(&request); err != nil {
		respondWithErr(http.StatusBadRequest, err, "invalid request", c)
		return
	}

	passwordHash, err := util.GeneratePassword(request.Password)
	if err != nil {
		respondWithErr(http.StatusInternalServerError, err, "error hashing password", c)
		return
	}

	createUserParams := db.CreateUserParams{
		Username: request.Username,
		Password: passwordHash,
		Email:    request.Email,
		FullName: request.FullName,
	}

	newUser, err := server.store.CreateUser(c, createUserParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Constraint {
			case "users_pkey":
				respondWithErr(http.StatusBadRequest, err, fmt.Sprintf("user with username %v exists", request.Username), c)
				return
			case "users_email_key":
				respondWithErr(http.StatusBadRequest, err, fmt.Sprintf("user with email %v exists", request.Email), c)
				return
			}
		}
		respondWithErr(http.StatusInternalServerError, err, "error creating user", c)
		return
	}

	token, err := server.tokenMaker.CreateToken(newUser.Username, server.config.TokenDuration)
	if err != nil {
		respondWithErr(http.StatusInternalServerError, fmt.Errorf("internal"), "error signing up user", c)
		return
	}

	cookieDuration := time.Now().Add(server.config.TokenDuration)
	c.SetCookie(auth_cookie_key, token, int(cookieDuration.Unix()), "/", "localhost", true, true)

	response := newUserResponse(newUser)
	c.JSON(http.StatusOK, response)

}
