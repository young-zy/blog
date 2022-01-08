package databases

import (
	"context"

	"github.com/young-zy/blog/models"
)

// AddPost adds a post to the database
func (tx *Transaction) AddPost(c context.Context, post *models.Post) error {
	return tx.WithContext(c).
		Create(post).
		Error
}

// GetPosts returns a list of posts, total count, and the error
func (tx *Transaction) GetPosts(c context.Context, page int, size int) (postList []*models.Post, totalCount int64, err error) {
	postDb := tx.WithContext(c).
		Model(&models.Post{})
	err = postDb.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&postList).
		Error
	return
}

// GetPost returns a post with the postID provided
func (tx *Transaction) GetPost(c context.Context, postID *uint) (post *models.Post, err error) {
	err = tx.WithContext(c).
		Model(&models.Post{}).
		Where("id = ?", *postID).
		Find(&post).
		Error
	return
}

// UpdatePost updates a post and saves to database
func (tx *Transaction) UpdatePost(c context.Context, post *models.Post) (rowsAffected int64, err error) {
	result := tx.WithContext(c).
		Where("id = ?", post.ID).
		Updates(post)
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

// DeletePost deletes a post with the postID provided
func (tx *Transaction) DeletePost(c context.Context, postID *uint) (rowsAffected int64, err error) {
	// delete all the replies of the post
	result := tx.WithContext(c).
		Where("posts_id = ?", *postID).
		Delete(&models.Reply{}).
		Where("id = ?", postID).
		Delete(&models.Post{})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}
