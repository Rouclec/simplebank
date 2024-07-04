package api

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rouclec/simplebank/db/sqlc"
	"github.com/rouclec/simplebank/token"
	"github.com/rouclec/simplebank/util"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	config     util.Config
}

// Creates a new HTTP server instance and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) //Switch between NewPasetoMaker and NewJWTMaker to use either Paseto or JWT tokens respectively
	if err != nil {
		return nil, fmt.Errorf("error creating token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
		v.RegisterValidation("email", validateEmail)
		v.RegisterValidation("password", validatePassword)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	//add routes to router
	router.POST("api/v1/auth/signup", server.createUser)
	router.POST("api/v1/auth/login", server.login)

	authRoutes := router.Group("/api/v1").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PATCH("/accounts", server.addAccountBalance)

	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.GET("/transfers/:id", server.getTransfer)
	authRoutes.GET("/transfers", server.listTransfers)

	authRoutes.GET("/users/:username", server.getUser)

	server.router = router
}

// start the HTTP server on the given address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"message": strings.Split(err.Error(), "\n")[0],
	}
}
