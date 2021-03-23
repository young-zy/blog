package databases

import (
	"context"

	"blog/models"
)

func (tx *Transaction) AddReply(c context.Context, reply *models.Reply) error {
	return tx.tx.WithContext(c).Create(reply).Error
}

func (tx *Transaction) GetReplies(c context.Context, postId *uint, page int, size int) (replyList []*models.Reply, totalCount int64, err error) {
	replyDB := tx.tx.WithContext(c).Model(&models.Reply{}).Where(&models.Reply{PostsId: *postId})
	err = replyDB.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&replyList).
		Error
	return
}

func (tx *Transaction) DeleteReply(c context.Context, replyId *uint) (rowsAffected int64, err error) {
	result := tx.tx.WithContext(c).Delete(&models.Reply{Id: replyId})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

func (tx *Transaction) UpdateReply(c context.Context, reply *models.Reply) error {
	return tx.tx.WithContext(c).
		Model(&models.Reply{}).
		Where("id = ?", *reply.Id).
		Updates(reply).Error
}

func (tx *Transaction) GetReply(c context.Context, replyId *uint) (reply *models.Reply, err error) {
	err = tx.tx.WithContext(c).
		Where("id = ?", *replyId).
		Find(reply).
		Error
	return
}
