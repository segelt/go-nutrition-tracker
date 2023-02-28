package middleware

import (
	e "nutritiontracker/resource/common"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err

		switch v := err.(type) {
		case e.AppError:
			c.JSON(v.Code, v)
		case error:
		default:
			c.JSON(-1, err.Error())
		}
	}
}
