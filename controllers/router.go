package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/middleware"
)

// Router the main router of the project
var Router = gin.Default()

func init() {
	Router.Use(cors())
	Router.Use(errorHandling())
	Router.Use(logging())
	Router.Use(middleware.APIRestrict.RestrictionMiddleware())
	Router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	Router.POST("/login", middleware.AuthMiddleware.LoginHandler)
	Router.POST("/logout", middleware.AuthMiddleware.LogoutHandler)
	Router.GET("/auth/token", middleware.AuthMiddleware.RefreshHandler)
	initQuestionGroup()
	initUserGroup()
	initPostGroup()
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, captchaToken")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func logging() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Next()
		errs := context.Errors.ByType(gin.ErrorTypePrivate)
		for _, err := range errs {
			trace, _ := context.Get("traceID")
			log.Printf("[trace-%s] Internal Error: %v", trace, err)
		}
	}
}

func errorHandling() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("traceID", uuid.New().String())
		context.Next()
		// error type, error meta, error string
		// _ = context.Error(errors.New("")).SetType(gin.ErrorTypePublic)
		errs := context.Errors.ByType(gin.ErrorTypePublic)
		errs = append(errs, context.Errors.ByType(gin.ErrorTypeBind)...)
		if len(errs) > 0 {
			err := errs.Last()
			switch err.Type {
			case gin.ErrorTypeBind:
				context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": fmt.Sprintf("failed to bind, %s", err.Err.Error()),
				})
			case gin.ErrorTypePublic:
				httpError := err.Err.(common.HTTPError)
				context.AbortWithStatusJSON(httpError.Status(), gin.H{
					"code":    httpError.Status(),
					"message": httpError.Message(),
				})
			default:
				context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": err.Err.Error(),
				})
			}
		}
	}
}
