package models

import "time"

// Question is used for base and orm
type Question struct {
	*QuestionResponse
	Email *string `gorm:"type:VARCHAR(100);" json:"email"`
}

// NewQuestionRequest is used when creating a new request
type NewQuestionRequest struct {
	QuestionContent string  `json:"questionContent" binding:"required"`
	Email           *string `json:"email" binding:"omitempty,email"`
}

// QuestionListResponse is used for list response
type QuestionListResponse struct {
	Questions  []*QuestionInListResponse `json:"questions"`
	TotalCount int64                     `json:"totalCount"`
}

// QuestionInListResponse is used in list response
type QuestionInListResponse struct {
	ID              *uint      `gorm:"primaryKey;type:INT;NOT NULL" json:"id"`
	QuestionContent string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"questionContent"`
	CreateTime      *time.Time `gorm:"type:DATETIME;" json:"createTime"`
	AnswerTime      *time.Time `gorm:"type:VARCHAR(100);" json:"answerTime"`
	IsAnswered      bool       `gorm:"type:TINYINT" json:"isAnswered"`
}

// QuestionResponse is used for single response
type QuestionResponse struct {
	*QuestionInListResponse
	AnswerContent *string `gorm:"type:MEDIUMTEXT;" json:"answerContent"`
}
