package models

import (
	"time"
)

// Post is used for base orm
type Post struct {
	ID          *uint     `gorm:"type:INT;NOT NULL"`
	Title       string    `gorm:"type:VARCHAR(100);NOT NULL"`
	Content     string    `gorm:"type:MEDIUMTEXT;NOT NULL"`
	Author      uint      `gorm:"type:INT;NOT NULL"`
	LastUpdated time.Time `gorm:"type:DATETIME;"`
}

// PostRequest is used in request
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// PostListResponse is used in list response
type PostListResponse struct {
	Posts      []*Post
	TotalCount int64
}

// PostResponse is used for single object response
type PostResponse struct {
	*Post
	Author *SimpleUser
}
