package handler

import (
	"context"
	"net/http"
	"nutritiontracker/handler/middleware"
	"nutritiontracker/resource"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
	router *gin.Engine

	MealService resource.MealService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: gin.Default(),
	}

	s.router.Use(middleware.ErrorHandler)
	mealGroup := s.router.Group("/meal")
	s.registerMealRoutes(mealGroup)

	return s
}

func (s *Server) Start() error {
	return s.router.Run("localhost:5820")
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
