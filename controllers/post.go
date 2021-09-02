package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/middleware"
	"github.com/young-zy/blog/models"
	"github.com/young-zy/blog/services"
)

func initPostGroup() {
	postGroup := Router.Group("/post")
	{
		postGroup.GET("", getPostList)
		postGroup.POST("", middleware.AuthMiddleware.MiddlewareFunc(), addPost)
		postGroup.GET("/:postID", getPost)
		postGroup.PUT("/:postID", middleware.AuthMiddleware.MiddlewareFunc())
		postGroup.DELETE("/:postID", middleware.AuthMiddleware.MiddlewareFunc(), deletePost)
	}
	replyGroup := postGroup.Group("/:postID/reply")
	{
		replyGroup.GET("", getReply)
		replyGroup.POST("", replyPost)
		replyGroup.PUT("/:replyID", middleware.AuthMiddleware.MiddlewareFunc(), updateReply)
		replyGroup.DELETE("/:replyID", middleware.AuthMiddleware.MiddlewareFunc(), deleteReply)
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
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("wrong format of postID provided")).SetType(gin.ErrorTypePublic)
		return
	}
	postResponse, ok := services.GetPost(c, common.IntToUintPointer(postID))
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
	services.AddPost(c, postRequest)
}

func deletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("wrong format of postID provided")).SetType(gin.ErrorTypePublic)
		return
	}
	services.DeletePost(c, common.IntToUintPointer(postID))
}

func getReply(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse postID")).SetType(gin.ErrorTypePublic)
	}
	pager := common.NewPager()
	if c.BindQuery(pager) != nil {
		return
	}
	services.GetReplies(c, common.IntToUintPointer(postID), pager.Page, pager.Size)
}

func replyPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse postID")).SetType(gin.ErrorTypePublic)
	}
	replyRequest := &models.ReplyRequest{}
	err = c.BindJSON(replyRequest)
	if err != nil {
		return
	}
	services.ReplyPost(c, replyRequest, common.IntToUintPointer(postID))
}

func updateReply(c *gin.Context) {
	replyID, err := strconv.Atoi(c.Param("replyID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse replyID"))
	}
	replyRequest := &models.ReplyRequest{}
	err = c.BindJSON(replyRequest)
	if err != nil {
		return
	}
	services.UpdateReply(c, replyRequest, common.IntToUintPointer(replyID))
}

func deleteReply(c *gin.Context) {
	replyID, err := strconv.Atoi(c.Param("replyID"))
	if err != nil {
		_ = c.Error(common.NewBadRequestError("failed to parse replyID"))
	}
	services.DeleteReply(c, common.IntToUintPointer(replyID))
}
