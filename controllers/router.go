package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog/common"
	"blog/middleware"
)

var Router = gin.Default()

func init() {
	Router.Use(Cors())
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	Router.POST("/login", middleware.AuthMiddleware.LoginHandler)
	Router.POST("/logout", middleware.AuthMiddleware.LogoutHandler)
	Router.GET("/auth/token", middleware.AuthMiddleware.RefreshHandler)
}

func handleError(c *gin.Context, httpError common.HttpError) {
	c.AbortWithStatusJSON(httpError.Status(), gin.H{
		"code":    httpError.Status(),
		"message": httpError.Message(),
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}
