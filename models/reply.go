package models

import "time"

// Reply is used for base and orm
type Reply struct {
	ID          *uint      `gorm:"type:INT;NOT NULL" json:"id"`
	Content     string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"content"`
	Email       string     `gorm:"type:VARCHAR(45);" json:"email"`
	PostsID     uint       `gorm:"type:INT;NOT NULL" json:"postsID"`
	LastUpdated *time.Time `gorm:"type:DATETIME;" json:"lastUpdated"`
}

// ReplyRequest is used on create or update request
type ReplyRequest struct {
	Content string `json:"content"`
	Email   string `json:"email"`
}

// ReplyResponse is used for response
type ReplyResponse struct {
	Replies    []*Reply
	TotalCount int64
}
