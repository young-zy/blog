package models

import (
	"time"
)

type Post struct {
	Id          *uint     `gorm:"type:INT;NOT NULL"`
	Title       string    `gorm:"type:VARCHAR(100);NOT NULL"`
	Content     string    `gorm:"type:MEDIUMTEXT;NOT NULL"`
	Author      uint      `gorm:"type:INT;NOT NULL"`
	LastUpdated time.Time `gorm:"type:DATETIME;"`
}

type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// struct used in response
type PostListResponse struct {
	Posts      []*Post
	TotalCount int64
}

type PostResponse struct {
	*Post
	Author *SimpleUser
}
