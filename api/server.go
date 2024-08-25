package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/techschool/simplebank/db/sqlc"
	"github.com/techschool/simplebank/token"
	"github.com/techschool/simplebank/util"
)

type Server struct {
	tokenMaker token.Maker
	config util.Config
	store db.Store
	router *gin.Engine
}

func NewServer (config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config: config,
		store: store,
		tokenMaker: tokenMaker,
	}
	

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/medium", server.createMedia)
	router.GET("/medium/:id", server.getMedia)
	router.GET("/medium", server.listMedia)
	router.DELETE("/medium/:id", server.deleteMedia)
	router.PATCH("/medium/:id", server.updateMedia)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/collections", server.createCollection)
	authRoutes.GET("/collections/:id", server.getCollection)
	authRoutes.GET("/collections", server.listCollections)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}