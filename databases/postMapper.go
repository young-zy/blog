package databases

import "blog/models"

// add a post to the database
func AddPost(post *models.Post) error {
	return db.Create(post).Error
}

// returns a list of posts, total count, and the error
func GetPosts(page int, size int) (postList []*models.Post, totalCount int64, err error) {
	postDb := db.Model(&models.Post{})
	err = postDb.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&postList).
		Error
	return
}

func UpdatePost(post *models.Post) (rowsAffected int64, err error) {
	// TODO wrong use of save
	result := db.Save(post)
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

func DeletePost(postId uint) (rowsAffected int64, err error) {
	// delete all the replies of the post
	result := db.Delete(&models.Reply{
		PostsId: postId,
		// delete the post
	}).Delete(&models.Post{
		Id: postId,
	})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}
