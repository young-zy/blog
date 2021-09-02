package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/middleware"
	"github.com/young-zy/blog/models"
	"github.com/young-zy/blog/services"
)

func initUserGroup() {
	userGroup := Router.Group("/user")
	{
		userGroup.GET("", middleware.AuthMiddleware.MiddlewareFunc(), getSelf)
		userGroup.POST("/avatar", middleware.AuthMiddleware.MiddlewareFunc(), setAvatar)
		userGroup.POST("", register)
		userGroup.GET("/:username", getUser)
	}
}

func getUser(c *gin.Context) {
	username := c.Param("username")
	user, ok := services.GetUser(c, username)
	if ok {
		c.JSON(http.StatusOK, user)
	}
}

func setAvatar(c *gin.Context) {
	req := &setAvatarRequest{}
	user, exists := c.Get("User")
	if !exists {
		_ = common.NewNotFoundError("user not found")
	}
	services.SetAvatar(c, user.(models.User).Username, req.Avatar)
}

type setAvatarRequest struct {
	// base 64 of avatar image
	Avatar string `json:"avatar" binding:"required"`
}

func getSelf(c *gin.Context) {
	user, ok := c.Get("User")
	if !ok {
		common.NewInternalError(c, errors.New("identity not found in context"))
	}
	c.JSON(http.StatusOK, user)
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
