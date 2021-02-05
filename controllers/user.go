package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog/models"
	"blog/services"
)

func init() {
	userGroup := Router.Group("/user")
	{
		userGroup.POST("", register)
	}
}

func getUser(c *gin.Context) {

}

func register(c *gin.Context) {
	body := &models.UserRegister{}
	err := c.BindJSON(body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "failed to parse body",
		})
	}
	httpError := services.Register(body.Username, body.Password, body.Email)
	if httpError != nil {
		// handle error
		handleError(c, httpError)
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"code":    http.StatusNoContent,
			"message": "user registered successfully",
		})
	}
}
