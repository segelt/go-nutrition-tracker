package handler

import (
	"fmt"
	"net/http"
	"nutritiontracker/resource"
	e "nutritiontracker/resource/common"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerAuthRoutes(r *gin.RouterGroup) {
	r.POST("/", s.Register)
}

func (s *Server) Register(c *gin.Context) {
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
