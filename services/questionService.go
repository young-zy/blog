package services

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"blog/common"
	"blog/databases"
	"blog/models"
)

func GetQuestion(ctx *gin.Context, questionId *uint) (question *models.QuestionResponse, ok bool) {
	ok = true
	question, err := databases.GetQuestion(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(common.NewNotFoundError("question not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(ctx, err)
		}
		ok = false
	}
	return
}

func GetQuestions(ctx *gin.Context, page int, size int, filter string) (questionListResponse *models.QuestionListResponse, ok bool) {
	ok = true
	questionList, totalCount, err := databases.GetQuestions(ctx, page, size, filter)
	if err != nil {
		// check if err is mysql error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(common.NewNotFoundError("no question available")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(ctx, err)
		}
		ok = false
		return
	}
	questionListResponse = &models.QuestionListResponse{
		Questions:  questionList,
		TotalCount: totalCount,
	}
	return
}

// add a question using databases.AddQuestion, returns if operation is successful
func AddQuestion(ctx *gin.Context, question *models.NewQuestionRequest) bool {
	err := databases.AddQuestion(ctx, question)
	if err != nil {
		common.NewInternalError(ctx, err)
		return false
	}
	return true
}

// returns if operation is successful
func AnswerQuestion(ctx *gin.Context, questionId *uint, content *string) bool {
	tx := databases.GetTransaction()
	question, err := databases.GetQuestionWithTransaction(ctx, tx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(common.NewNotFoundError("question not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(ctx, err)
		}
		tx.Rollback()
		return false
	}
	if question.AnswerContent != nil && *question.AnswerContent != "" {
		tx.Rollback()
		_ = ctx.Error(common.NewBadRequestError("question has already been answered")).SetType(gin.ErrorTypePublic)
		return false
	}
	question.AnswerContent = content
	question.IsAnswered = true
	err = databases.UpdateQuestionWithTransaction(ctx, tx, question)
	if err != nil {
		tx.Rollback()
		common.NewInternalError(ctx, err)
		return false
	}
	tx.Commit()
	// notify the email
	if question.Email != nil && *question.Email != "" {
		link := fmt.Sprintf("https://young-zy.com/question/%d", *questionId)
		message := fmt.Sprintf("您的提问已被回复: %s", link)
		title := "您在提问箱的提问有新回答"
		go common.SendMail(ctx.Copy(), *question.Email, title, message)
	}
	return true
}

func UpdateAnswer(ctx *gin.Context, questionId *uint, content *string) bool {
	tx := databases.GetTransaction()
	question, err := databases.GetQuestionWithTransaction(ctx, tx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(common.NewNotFoundError("question not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(ctx, err)
		}
		tx.Rollback()
		return false
	}
	question.AnswerContent = content
	err = databases.UpdateQuestionWithTransaction(ctx, tx, question)
	if err != nil {
		common.NewInternalError(ctx, err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true
}

func DeleteQuestion(ctx *gin.Context, questionId *uint) bool {
	err := databases.DeleteQuestion(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = ctx.Error(common.NewNotFoundError("question not found")).SetType(gin.ErrorTypePublic)
		} else {
			common.NewInternalError(ctx, err)
		}
		return false
	}
	return true
}
