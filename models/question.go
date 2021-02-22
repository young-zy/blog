package models

import "time"

type Question struct {
	Id              *uint      `gorm:"primaryKey;type:INT;NOT NULL" json:"id"`
	QuestionContent string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"questionContent"`
	CreateTime      *time.Time `gorm:"type:DATETIME;" json:"createTime"`
	Email           *string    `gorm:"type:VARCHAR(100);" json:"email"`
	AnswerContent   *string    `gorm:"type:MEDIUMTEXT;" json:"answerContent"`
	AnswerTime      *time.Time `gorm:"type:DATETIME;" json:"answerTime"`
	IsAnswered      bool       `gorm:"type:TINYINT" json:"isAnswered"`
}

type NewQuestionRequest struct {
	QuestionContent string  `json:"questionContent" binding:"required"`
	Email           *string `json:"email" binding:"omitempty,email"`
}

type QuestionListResponse struct {
	Questions  []*QuestionInListResponse `json:"questions"`
	TotalCount int64                     `json:"totalCount"`
}

type QuestionInListResponse struct {
	Id              uint       `json:"id"`
	QuestionContent string     `json:"questionContent"`
	CreateTime      *time.Time `json:"createTime"`
	AnswerTime      *time.Time `json:"answerTime"`
	IsAnswered      bool       `json:"isAnswered"`
}

type QuestionResponse struct {
	*QuestionInListResponse
	AnswerContent string `json:"answerContent"`
}
