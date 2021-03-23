package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog/common"
	"blog/middleware"
	"blog/models"
	"blog/services"
)

func initPostGroup() {
	postGroup := Router.Group("/post")
	{
		postGroup.GET("", getPostList)
		postGroup.POST("", middleware.AuthMiddleware.MiddlewareFunc(), addPost)
		postGroup.GET("/:postId", getPost)
		postGroup.PUT("/:postId", middleware.AuthMiddleware.MiddlewareFunc())
		postGroup.DELETE("/:postId", middleware.AuthMiddleware.MiddlewareFunc(), deletePost)
	}
	replyGroup := postGroup.Group("/:postId/reply")
	{
		replyGroup.GET("", getReply)
		replyGroup.POST("", replyPost)
		replyGroup.PUT("/:replyId", middleware.AuthMiddleware.MiddlewareFunc(), updateReply)
		replyGroup.DELETE("/:replyId", middleware.AuthMiddleware.MiddlewareFunc(), deleteReply)
	}
}

func getPostList(c *gin.Context) {
	pager := common.NewPager()
	if c.BindQuery(pager) != nil {
		return
	}
	postListResponse, ok := services.GetAllPosts(c, pager.Page, pager.Size)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, postListResponse)
}

func getPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("wrong format of postId provided")).SetType(gin.ErrorTypePublic)
		return
	}
	postResponse, ok := services.GetPost(c, common.IntToUintPointer(postId))
	if !ok {
		return
	}
	c.JSON(http.StatusOK, postResponse)
}

func addPost(c *gin.Context) {
	postRequest := &models.PostRequest{}
	err := c.BindJSON(postRequest)
	if err != nil {
		return
	}
	ok := services.AddPost(c, postRequest)
	if !ok {
		return
	}
}

func deletePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("wrong format of postId provided")).SetType(gin.ErrorTypePublic)
		return
	}
	services.DeletePost(c, common.IntToUintPointer(postId))
}

func getReply(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse postId")).SetType(gin.ErrorTypePublic)
	}
	pager := common.NewPager()
	if c.BindQuery(pager) != nil {
		return
	}
	services.GetReplies(c, common.IntToUintPointer(postId), pager.Page, pager.Size)
}

func replyPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse postId")).SetType(gin.ErrorTypePublic)
	}
	replyRequest := &models.ReplyRequest{}
	err = c.BindJSON(replyRequest)
	if err != nil {
		return
	}
	services.ReplyPost(c, replyRequest, common.IntToUintPointer(postId))
}

func updateReply(c *gin.Context) {
	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse replyId"))
	}
	replyRequest := &models.ReplyRequest{}
	err = c.BindJSON(replyRequest)
	if err != nil {
		return
	}
	services.UpdateReply(c, replyRequest, common.IntToUintPointer(replyId))
}

func deleteReply(c *gin.Context) {
	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse replyId"))
	}
	services.DeleteReply(c, common.IntToUintPointer(replyId))
}
