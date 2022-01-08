package databases

import (
	"context"

	"github.com/young-zy/blog/models"
)

// AddReply adds a reply to a post
func (tx Transaction) AddReply(c context.Context, reply *models.Reply) error {
	return tx.WithContext(c).Create(reply).Error
}

// GetReplies gets a list of replies of a post
func (tx Transaction) GetReplies(c context.Context, postID *uint, page int, size int) (replyList []*models.Reply, totalCount int64, err error) {
	replyDB := tx.WithContext(c).Model(&models.Reply{}).Where(&models.Reply{PostsID: *postID})
	err = replyDB.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Find(&replyList).
		Error
	return
}

// DeleteReply deletes reply with the ID provided
func (tx Transaction) DeleteReply(c context.Context, replyID *uint) (rowsAffected int64, err error) {
	result := tx.WithContext(c).Delete(&models.Reply{ID: replyID})
	rowsAffected = result.RowsAffected
	err = result.Error
	return
}

// UpdateReply updates a reply
func (tx Transaction) UpdateReply(c context.Context, reply *models.Reply) error {
	return tx.WithContext(c).
		Model(&models.Reply{}).
		Where("id = ?", *reply.ID).
		Updates(reply).Error
}

// GetReply gets a reply with the ID provided
func (tx Transaction) GetReply(c context.Context, replyID *uint) (reply *models.Reply, err error) {
	err = tx.WithContext(c).
		Where("id = ?", *replyID).
		Find(reply).
		Error
	return
}
