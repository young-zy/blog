package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog/models"
	"blog/services"
)

func initUserGroup() {
	userGroup := Router.Group("/user")
	{
		userGroup.POST("", register)
		userGroup.GET("/:username", getUser)
	}
}

func getUser(c *gin.Context) {
	username := c.Param("userId")
	services.GetUser(c, username)
}

func register(c *gin.Context) {
	body := &models.UserRegister{}
	err := c.BindJSON(body)
	if err != nil {
		return
	}
	if services.Register(c, body.Username, body.Password, body.Email) {
		c.JSON(http.StatusNoContent, gin.H{
			"code":    http.StatusNoContent,
			"message": "user registered successfully",
		})
	}
}
