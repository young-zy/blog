package models

import "time"

type Reply struct {
	Id          uint       `gorm:"type:INT;NOT NULL" json:"id"`
	Content     string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"content"`
	Email       string     `gorm:"type:VARCHAR(45);" json:"email"`
	PostsId     uint       `gorm:"type:INT;NOT NULL" json:"postsId"`
	LastUpdated *time.Time `gorm:"type:DATETIME;" json:"lastUpdated"`
}
