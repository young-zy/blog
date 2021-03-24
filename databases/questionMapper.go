package databases

import (
	"context"

	"gorm.io/gorm"

	"blog/common"
	"blog/models"
)

// GetQuestions Get a list of questions, filter could be 'solved' or 'unsolved', any other values(including unprovided) will be ignored.
func (tx *Transaction) GetQuestions(ctx context.Context, page int, size int, filter string) (questionList []*models.QuestionInListResponse, totalCount int64, err error) {
	questionDb := tx.tx.WithContext(ctx).
		Table("questions").
		Model(&models.QuestionResponse{})
	if filter == "solved" {
		questionDb = questionDb.Where("is_answered = 1")
	} else if filter == "unsolved" {
		questionDb = questionDb.Where("is_answered = 0")
	}
	err = questionDb.
		Count(&totalCount).
		Offset((page - 1) * size).
		Limit(size).
		Order("is_answered desc").
		Order("id desc").
		Find(&questionList).
		Error
	return
}

//// get a list of questions in transaction
//func GetQuestionsWithTransaction(ctx context.Context, tx *gorm.DB, page int, size int) (questionList []*models.QuestionResponse, totalCount int64, err error) {
//	questionDb := tx.WithContext(ctx).Table("questions").Model(&models.QuestionResponse{})
//	err = questionDb.
//		Count(&totalCount).
//		Offset((page - 1) * size).
//		Limit(size).
//		Order("is_answered desc").
//		Order("id desc").
//		Find(&questionList).
//		Error
//	return
//}

//// get a single question bu questionID
//func GetQuestion(ctx context.Context, questionID *uint) (question *models.QuestionResponse, err error) {
//	question = &models.QuestionResponse{}
//	err = DefaultDb.WithContext(ctx).Model(&models.Question{}).Where(&models.Question{ID: questionID}).First(question).Error
//	return
//}

// GetQuestion get a single question with the ID provided
func (tx *Transaction) GetQuestion(ctx context.Context, questionID *uint) (question *models.Question, err error) {
	question = &models.Question{}
	err = tx.tx.WithContext(ctx).
		Where("id = ?", questionID).
		First(question).
		Error
	return
}

// AddQuestion add a question to database
func (tx *Transaction) AddQuestion(ctx context.Context, question *models.NewQuestionRequest) error {
	return tx.tx.WithContext(ctx).Create(&models.Question{
		QuestionResponse: &models.QuestionResponse{
			QuestionInListResponse: &models.QuestionInListResponse{
				QuestionContent: question.QuestionContent,
				CreateTime:      common.Now(),
				IsAnswered:      false,
			},
		},
		Email: question.Email,
	}).Error
}

//// update a question object, must check existence in a transaction before calling this method
//func UpdateQuestion(ctx context.Context, question *models.Question) error {
//	return DefaultDb.WithContext(ctx).Model(&models.Question{}).Updates(question).Error
//}

// UpdateQuestion update a question object
func (tx *Transaction) UpdateQuestion(ctx context.Context, question *models.Question) error {
	question.AnswerTime = common.Now()
	return tx.tx.Session(&gorm.Session{AllowGlobalUpdate: true}).
		WithContext(ctx).
		Model(&models.Question{}).
		Where("id = ?", question.ID).
		Updates(question).
		Error
}

// DeleteQuestion deletes a question with the ID provided
func (tx *Transaction) DeleteQuestion(ctx context.Context, questionID *uint) error {
	return tx.tx.WithContext(ctx).
		Delete(&models.Question{}, questionID).
		Error
}
