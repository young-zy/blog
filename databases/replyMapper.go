package databases

import "blog/models"

func AddReply(reply *models.Reply) error {
	return db.Create(reply).Error
}

func GetReplies(postId uint, page int, size int) (replyList []*models.Reply, totalCount int64, err error) {
	replyDB := db.Model(&models.Reply{}).Where(&models.Reply{PostsId: postId})
	err = replyDB.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&replyList).
		Error
	return
}

func DeleteReply(replyId uint) (rowsAffected int64, err error) {
	result := db.Delete(&models.Reply{Id: replyId})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

func UpdateReply(reply *models.Reply) {
	db.Save(reply)
}
