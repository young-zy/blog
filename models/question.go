package models

import "time"

type Question struct {
	*QuestionResponse
	Email *string `gorm:"type:VARCHAR(100);" json:"email"`
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
	Id              *uint      `gorm:"primaryKey;type:INT;NOT NULL" json:"id"`
	QuestionContent string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"questionContent"`
	CreateTime      *time.Time `gorm:"type:DATETIME;" json:"createTime"`
	AnswerTime      *time.Time `gorm:"type:VARCHAR(100);" json:"answerTime"`
	IsAnswered      bool       `gorm:"type:TINYINT" json:"isAnswered"`
}

type QuestionResponse struct {
	*QuestionInListResponse
	AnswerContent *string `gorm:"type:MEDIUMTEXT;" json:"answerContent"`
}
