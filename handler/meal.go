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

type MealServer struct {
	BaseServer  *BaseServer
	MealService resource.MealService
}

func NewMealServer(serverConf configs.ServerConfig) *MealServer {
	baseServer := &BaseServer{
		server:     &http.Server{},
		router:     gin.Default(),
		ServerConf: serverConf,
	}
	s := &MealServer{
		BaseServer: baseServer,
	}

	baseServer.router.Use(middleware.ErrorHandler)

	return s
}

func (s *MealServer) Start() error {
	s.registerRoutes()
	return s.BaseServer.Start()
}

func (s *MealServer) Close() error {
	return s.BaseServer.Close()
}

func (s *MealServer) registerRoutes() {
	r := s.BaseServer.router.Group("/meal")
	authMiddleware := middleware.NewAuthMiddleware(s.BaseServer.ServerConf.JWT_SECRET)

	r.Use(authMiddleware.WithAuthentication())
	{
		r.POST("/", s.Create)
		r.GET("/single", s.ById)
		r.GET("/filter", s.ListMeals)
		r.POST("/update", s.Update)
		r.DELETE("/", s.DeleteById)
	}
}

func (s *MealServer) Create(c *gin.Context) {
	var req resource.MealCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Create Meal input params | %s", err.Error()),
		})
		return
	}

	err := s.MealService.CreateMeal(&req)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Create Meal unsuccesful | %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *MealServer) ById(c *gin.Context) {
	targetId, ok := c.GetQuery("id")
	if !ok {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid or missing id input for find meal"),
		})
		return
	}

	targetMeal, err := s.MealService.FindMealById(&targetId)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Find Meal unsuccesful | %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, targetMeal)
}

func (s *MealServer) ListMeals(c *gin.Context) {
	var req resource.MealFilter
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("List Meal input params | %s", err.Error()),
		})
		return
	}

	meals, err := s.MealService.ListMeals(req)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Create Meal unsuccesful | %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, meals)
}

func (s *MealServer) Update(c *gin.Context) {
	var req resource.MealUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Update Meal input params | %s", err.Error()),
		})
		return
	}

	err := s.MealService.UpdateMeal(&req)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Update Meal unsuccesful | %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}

func (s *MealServer) DeleteById(c *gin.Context) {
	targetId, ok := c.GetQuery("id")
	if !ok {
		c.Error(e.AppError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("Invalid or missing id input for delete meal"),
		})
		return
	}

	err := s.MealService.DeleteMeal(targetId)
	if err != nil {
		c.Error(e.AppError{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Delete Meal unsuccesful | %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, "")
}
