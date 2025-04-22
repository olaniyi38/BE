package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/olaniyi38/BE/db/sqlc"
	"github.com/olaniyi38/BE/token"
	"github.com/olaniyi38/BE/util"
)

var auth_cookie_key = "auth_cookie"

// responsible for serving http requests for backend service
// this is what we attach endpoints to
type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("error creating token maker %v ", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	router := gin.Default()
	server.router = router

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("email", validateEmail)
	}

	// server.router.Use(server.authMiddleware)
	setupRouter(server)

	return server, nil
}

func setupRouter(server *Server) {
	protectedRoutes := server.router.Group("/").Use(authMiddleware(server.tokenMaker))

	//accounts
	protectedRoutes.POST("/createAccount", server.createAccount)
	protectedRoutes.GET("/account/:id", server.getAccount)
	protectedRoutes.GET("/accounts", server.listAccounts)

	//transfers
	protectedRoutes.POST("/createTransfer", server.createTransfer)

	//users
	protectedRoutes.GET("/user", server.getUser)

	//auth
	server.router.POST("/auth/signUp", server.SignUp)
	server.router.POST("auth/login", server.Login)
}

// start the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func respondWithErr(code int, err error, msg string, c *gin.Context) {
	c.JSON(code, gin.H{"error": err.Error(), "message": msg})
}
