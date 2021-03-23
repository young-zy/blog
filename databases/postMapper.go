package databases

import (
	"context"

	"blog/models"
)

// add a post to the database
func (tx *Transaction) AddPost(c context.Context, post *models.Post) error {
	return tx.tx.WithContext(c).Create(post).Error
}

// returns a list of posts, total count, and the error
func (tx *Transaction) GetPosts(c context.Context, page int, size int) (postList []*models.Post, totalCount int64, err error) {
	postDb := tx.tx.WithContext(c).Model(&models.Post{})
	err = postDb.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&postList).
		Error
	return
}

func (tx *Transaction) GetPost(c context.Context, postId *uint) (post *models.Post, err error) {
	err = tx.tx.WithContext(c).
		Model(&models.Post{}).
		Where("id = ?", *postId).
		Find(&post).
		Error
	return
}

func (tx *Transaction) UpdatePost(c context.Context, post *models.Post) (rowsAffected int64, err error) {
	result := tx.tx.WithContext(c).Where("id = ?", post.Id).Updates(post)
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

func (tx *Transaction) DeletePost(c context.Context, postId *uint) (rowsAffected int64, err error) {
	// delete all the replies of the post
	result := tx.tx.WithContext(c).
		Where("posts_id = ?", *postId).
		Delete(&models.Reply{}).
		Where("id = ?", postId).
		Delete(&models.Post{})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}
