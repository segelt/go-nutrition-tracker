package handler

import (
	"fmt"
	"net/http"
	"nutritiontracker/handler/middleware"
	"nutritiontracker/resource"
	e "nutritiontracker/resource/common"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerMealRoutes(r *gin.RouterGroup) {
	r.Use(middleware.WithAuthentication())
	{
		r.POST("/", s.Create)
		r.GET("/single", s.ById)
		r.GET("/filter", s.ListMeals)
		r.POST("/update", s.Update)
		r.DELETE("/", s.DeleteById)
	}
}

func (s *Server) Create(c *gin.Context) {
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

func (s *Server) ById(c *gin.Context) {
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

func (s *Server) ListMeals(c *gin.Context) {
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

func (s *Server) Update(c *gin.Context) {
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

func (s *Server) DeleteById(c *gin.Context) {
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
