package models

import (
	"time"
)

type Post struct {
	Id          uint      `gorm:"type:INT;NOT NULL"`
	Title       string    `gorm:"type:VARCHAR(100);NOT NULL"`
	Content     string    `gorm:"type:MEDIUMTEXT;NOT NULL"`
	Author      uint      `gorm:"type:INT;NOT NULL"`
	LastUpdated time.Time `gorm:"type:DATETIME;"`
}

// struct used in response
type PostResponse struct {
	Id          uint
	Title       string
	Content     string
	Author      uint
	LastUpdated *time.Time

	TotalCount  int64
	CurrentPage int
	Replies     []*Reply
}
