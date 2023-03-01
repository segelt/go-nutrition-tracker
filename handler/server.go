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
	AuthService resource.AuthService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{},
		router: gin.Default(),
	}

	s.router.Use(middleware.ErrorHandler)

	return s
}

func (s *Server) registerRoutes() {
	if s.MealService != nil {
		mealGroup := s.router.Group("/meal")
		s.registerMealRoutes(mealGroup)
	}

	if s.AuthService != nil {
		authGroup := s.router.Group("/auth")
		s.registerAuthRoutes(authGroup)
	}
}

func (s *Server) Start() error {
	s.registerRoutes()
	return s.router.Run("localhost:5820")
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
