package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog/common"
	"blog/databases"
	"blog/models"
)

// GetAllPosts returns the post list
func GetAllPosts(c *gin.Context, page int, size int) (resp *models.PostListResponse, ok bool) {
	posts, totCount, err := databases.Default.GetPosts(c, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("no post available")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		return
	}
	resp = &models.PostListResponse{
		Posts:      posts,
		TotalCount: totCount,
	}
	ok = true
	return
}

// GetPost returns a single post
func GetPost(c *gin.Context, postID *uint) (resp *models.PostResponse, ok bool) {
	post, err := databases.Default.GetPost(c, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("post not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		return
	}
	author, err := databases.Default.GetUserByID(c, &post.Author)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("post not found")).SetType(gin.ErrorTypePublic)
			author = &models.User{
				Username: "deleted user",
			}
		} else {
			common.NewInternalError(c, err)
			return
		}
	}
	resp = &models.PostResponse{
		Post:   post,
		Author: author.GetSimpleUser(),
	}
	//replies, replyCount, err := databases.Default.GetReplies(c, postID, page, size)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		resp.Replies = make([]*models.Reply, 0)
	//	} else {
	//		common.NewInternalError(c, err)
	//	}
	//	return
	//}
	//resp.Replies = replies
	//resp.ReplyCount = replyCount
	ok = true
	return
}

// AddPost creates a new post
func AddPost(c *gin.Context, post *models.PostRequest) (ok bool) {
	// get user from context
	userInterface, exists := c.Get("User")
	if !exists {
		common.NewInternalError(c, errors.New("user not found in context"))
	}
	user := userInterface.(*models.User)
	// enforce
	res, err := enforcer.Enforce(user, "", "add post")
	if err != nil {
		_ = c.Error(common.NewInternalServerError("failed to enforce")).SetType(gin.ErrorTypePrivate)
	}
	if !res {
		_ = c.Error(common.NewForbiddenError("permission denied")).SetType(gin.ErrorTypePublic)
		return
	}
	tx := databases.GetTransaction()
	err = tx.AddPost(c, &models.Post{
		Title:       post.Title,
		Content:     post.Content,
		Author:      *user.ID,
		LastUpdated: time.Time{},
	})
	if err != nil {
		common.NewInternalError(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	ok = true
	return
}

// DeletePost deletes a post
func DeletePost(c *gin.Context, postID *uint) (ok bool) {
	tx := databases.GetTransaction()
	rowsAffected, err := tx.DeletePost(c, postID)
	if err != nil {
		common.NewInternalError(c, err)
		tx.Rollback()
		return
	}
	if rowsAffected == 0 {
		_ = c.Error(common.NewNotFoundError("post not found")).SetType(gin.ErrorTypePublic)
		tx.Rollback()
		return
	}
	tx.Commit()
	ok = true
	return
}

// GetReplies returns a list of replies
func GetReplies(c *gin.Context, postID *uint, page int, size int) (replyResponse *models.ReplyResponse, ok bool) {
	replies, replyCount, err := databases.Default.GetReplies(c, postID, page, size)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			replyResponse = &models.ReplyResponse{
				Replies:    make([]*models.Reply, 0),
				TotalCount: replyCount,
			}
		} else {
			common.NewInternalError(c, err)
		}
		return
	}
	replyResponse = &models.ReplyResponse{
		Replies:    replies,
		TotalCount: replyCount,
	}
	ok = true
	return
}

// ReplPost replies a post
func ReplyPost(c *gin.Context, reply *models.ReplyRequest, postID *uint) (ok bool) {
	tx := databases.GetTransaction()
	_, err := tx.GetPost(c, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("post not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		tx.Rollback()
		return
	}
	now := time.Now()
	err = tx.AddReply(c, &models.Reply{
		Content:     reply.Content,
		Email:       reply.Email,
		PostsID:     *postID,
		LastUpdated: &now,
	})
	if err != nil {
		common.NewInternalError(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	ok = true
	return
}

// UpdateReply updates a reply, works only when a user is logged in and its email is the same with the reply
func UpdateReply(c *gin.Context, replyRequest *models.ReplyRequest, replyID *uint) (ok bool) {
	tx := databases.GetTransaction()
	reply, err := tx.GetReply(c, replyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(common.NewNotFoundError("reply not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(c, err)
		}
		tx.Rollback()
		return
	}
	// add enforcement policy and uncomment the code
	//user, exists := c.Get("User")
	//if !exists {
	//	common.NewInternalError(c, errors.New("user object not found in context"))
	//}
	//res, err := common.enforcer.Enforce(user, reply, "update reply")
	//if err != nil {
	//	_ = c.Error(errors.New("enforcer encountered an error")).SetType(gin.ErrorTypePrivate)
	//}
	//if !res {
	//	tx.Rollback()
	//	_ = c.Error(common.NewForbiddenError("permission denied")).SetType(gin.ErrorTypePublic)
	//	return
	//}
	reply.Content = replyRequest.Content
	reply.LastUpdated = common.Now()
	err = tx.UpdateReply(c, reply)
	if err != nil {
		common.NewInternalError(c, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	ok = true
	return
}

// DeleteReply deletes a reply
func DeleteReply(c *gin.Context, replyID *uint) (ok bool) {
	tx := databases.GetTransaction()
	rowsAffected, err := tx.DeleteReply(c, replyID)
	if err != nil {
		common.NewInternalError(c, err)
	}
	if rowsAffected == 0 {
		_ = c.Error(common.NewNotFoundError("reply not found"))
	}
	tx.Commit()
	ok = true
	return
}
