package models

import "time"

type Question struct {
	Id              uint       `gorm:"type:INT;NOT NULL" json:"id"`
	QuestionContent string     `gorm:"type:MEDIUMTEXT;NOT NULL" json:"question_content"`
	CreateTime      *time.Time `gorm:"type:DATETIME;" json:"create_time"`
	Email           string     `gorm:"type:VARCHAR(100);" json:"email"`
	AnswerContent   string     `gorm:"type:MEDIUMTEXT;" json:"answer_content"`
	AnswerTime      *time.Time `gorm:"type:DATETIME;" json:"answer_time"`
	IsAnswered      bool       `gorm:"type:TINYINT" json:"is_answered"`
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
