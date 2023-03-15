package handler

import (
	"fmt"
	"net/http"
	"nutritiontracker/configs"
	"nutritiontracker/handler/middleware"
	"nutritiontracker/resource"
	e "nutritiontracker/resource/common"

	"github.com/gin-gonic/gin"
)

type AuthServer struct {
	BaseServer  *BaseServer
	AuthService resource.AuthService
}

func NewAuthServer(serverConf configs.ServerConfig) *AuthServer {
	baseServer := &BaseServer{
		server:     &http.Server{},
		router:     gin.Default(),
		ServerConf: serverConf,
	}
	s := &AuthServer{
		BaseServer: baseServer,
	}

	baseServer.router.Use(middleware.ErrorHandler)

	return s
}

func (s *AuthServer) Start() error {
	s.registerRoutes()
	return s.BaseServer.Start()
}

func (s *AuthServer) Close() error {
	return s.BaseServer.Close()
}

func (s *AuthServer) registerRoutes() {
	r := s.BaseServer.router.Group("/auth")
	r.POST("/", s.Register)
	r.POST("/login", s.Login)
}

func (s *AuthServer) Register(c *gin.Context) {
	var req resource.UserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("User registration input params | %s", err.Error()),
		})
		return
	}

	err := s.AuthService.RegisterUser(req)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("User register unsuccesful | %s", err.Error()),
		})
		return
	}

	c.Status(http.StatusCreated)
}

func (s *AuthServer) Login(c *gin.Context) {
	var req resource.UserInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("User login input params | %s", err.Error()),
		})
		return
	}

	generatedToken, err := s.AuthService.LoginUser(req)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("User login unsuccesful | %s", err.Error()),
		})
		return
	}

	c.String(http.StatusOK, generatedToken)
}
