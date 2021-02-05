package databases

import (
	"context"
	"gorm.io/gorm"

	"blog/models"
)

// get a list of questions
func GetQuestions(ctx context.Context, page int, size int) (questionList []*models.QuestionResponse, totalCount int64, err error) {
	questionDb := db.WithContext(ctx).Table("questions").Model(&models.QuestionResponse{})
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

// get a list of questions in transaction
func GetQuestionsWithTransaction(ctx context.Context, tx *gorm.DB, page int, size int) (questionList []*models.QuestionResponse, totalCount int64, err error) {
	questionDb := tx.WithContext(ctx).Table("questions").Model(&models.QuestionResponse{})
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

// get a single question bu questionId
func GetQuestion(ctx context.Context, questionId uint) (question *models.Question, err error) {
	question = &models.Question{}
	err = db.WithContext(ctx).Where(&models.Question{Id: questionId}).First(question).Error
	return
}

func GetQuestionWithTransaction(ctx context.Context, tx *gorm.DB, questionId uint) (question *models.Question, err error) {
	question = &models.Question{}
	err = tx.WithContext(ctx).Where(&models.Question{Id: questionId}).First(question).Error
	return
}

// add a question to database
func AddQuestion(ctx context.Context, question *models.Question) error {
	return db.WithContext(ctx).Create(question).Error
}

// update a question object
func UpdateQuestion(ctx context.Context, question *models.Question) error {
	return db.WithContext(ctx).Model(&models.Question{}).Updates(question).Error
}

// update a question object
func UpdateQuestionWithTransaction(ctx context.Context, tx *gorm.DB, question *models.Question) error {
	return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).WithContext(ctx).Model(&models.Question{}).Updates(question).Error
}

func DeleteQuestion(ctx context.Context, questionId uint) error {
	return db.WithContext(ctx).Delete(&models.Question{Id: questionId}).Error
}
