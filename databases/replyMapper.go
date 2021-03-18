package databases

import "blog/models"

func (tx *Transaction) AddReply(reply *models.Reply) error {
	return tx.tx.Create(reply).Error
}

func (tx *Transaction) GetReplies(postId uint, page int, size int) (replyList []*models.Reply, totalCount int64, err error) {
	replyDB := tx.tx.Model(&models.Reply{}).Where(&models.Reply{PostsId: postId})
	err = replyDB.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&replyList).
		Error
	return
}

func (tx *Transaction) DeleteReply(replyId uint) (rowsAffected int64, err error) {
	result := tx.tx.Delete(&models.Reply{Id: replyId})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

func (tx *Transaction) UpdateReply(reply *models.Reply) {
	tx.tx.Save(reply)
}
