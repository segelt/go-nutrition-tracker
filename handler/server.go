package handler

import (
	"context"
	"fmt"
	"net/http"
	"nutritiontracker/configs"
	"time"

	"github.com/gin-gonic/gin"
)

type Server interface {
	RegisterRoutes()
	Start() error
	Close() error
}

type BaseServer struct {
	server *http.Server
	router *gin.Engine

	ServerConf configs.ServerConfig
}

func (s *BaseServer) Start() error {
	return s.router.Run(fmt.Sprintf("localhost:%d", s.ServerConf.Port))
}

// Close gracefully shuts down the server.
func (s *BaseServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
