package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/rouclec/simplebank/db/sqlc"
)

// Server serves HTTP requests for our banking service
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Creates a new HTTP server instance and setup routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//add routes to router
	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.GET("/api/v1/accounts", server.listAccounts)

	server.router = router
	return server
}

//start the HTTP server on the given address
func (server *Server) Start(address string) error{
	return server.router.Run(address)
}


func errorResponse(err error) gin.H{
	return gin.H{
		"message": err.Error(),
	}
}