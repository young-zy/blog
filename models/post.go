package models

import (
	"time"
)

// Post is used for base orm
type Post struct {
	ID          *uint     `gorm:"type:INT;NOT NULL" json:"id"`
	Title       string    `gorm:"type:VARCHAR(100);NOT NULL" json:"title"`
	Content     string    `gorm:"type:MEDIUMTEXT;NOT NULL" json:"content"`
	Image       string    `gorm:"type:TEXT" json:"image"`
	Author      *User     `gorm:"foreignKey:ID;references:AuthorId"`
	AuthorId    uint      `gorm:"column:author;type:INT;NOT NULL" json:"author"`
	LastUpdated time.Time `gorm:"type:DATETIME;" json:"lastUpdated"`
}

// PostRequest is used in request
type PostRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// PostListResponse is used in list response
type PostListResponse struct {
	Posts      []*Post `json:"posts"`
	TotalCount int64   `json:"totalCount"`
}

// PostResponse is used for single object response
type PostResponse struct {
	*Post
	Author *SimpleUser
}
