package services

import (
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"

	"blog/common"
	"blog/databases"
	"blog/models"
)

func GetQuestions(ctx context.Context, page int, size int) (questionListResponse *models.QuestionListResponse, httpError common.HttpError) {
	questionList, totalCount, err := databases.GetQuestions(ctx, page, size)
	if err != nil {
		httpError = common.NewInternalServerError("error retrieving questions")
		return
	}
	questionListResponse = &models.QuestionListResponse{
		Questions:  questionList,
		TotalCount: totalCount,
	}
	return
}

// add a question using databases.AddQuestion
func AddQuestion(ctx context.Context, question *models.Question) (httpError common.HttpError) {
	err := databases.AddQuestion(ctx, question)
	if err != nil {
		httpError = common.NewInternalServerError(err.Error())
	}
	return
}

func AnswerQuestion(ctx context.Context, questionId uint, content string) (httpError common.HttpError) {
	tx := databases.GetTransaction()
	question, err := databases.GetQuestionWithTransaction(ctx, tx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("question not found")
			return
		}
		httpError = common.NewInternalServerError(err.Error())
		tx.Rollback()
		return
	}
	if question.AnswerContent != "" {
		tx.Rollback()
		return common.NewBadRequestError("question has already been answered")
	}
	question.AnswerContent = content
	err = databases.UpdateQuestionWithTransaction(ctx, tx, question)
	if err != nil {
		log.Println(err)
		httpError = common.NewInternalServerError(err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	// notify the email
	if question.Email != "" {
		link := fmt.Sprintf("https://young-zy.com/question/%d", questionId)
		message := fmt.Sprintf("您的提问已被回复: %s", link)
		c, cancel := context.WithCancel(ctx)
		defer cancel()
		title := "您在提问箱的提问有新回答"
		go common.SendMail(c, question.Email, title, message)
	}
	return
}

func UpdateAnswer(ctx context.Context, questionId uint, content string) (httpError common.HttpError) {
	tx := databases.GetTransaction()
	question, err := databases.GetQuestionWithTransaction(ctx, tx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("question not found")
			return
		}
		httpError = common.NewInternalServerError(err.Error())
		tx.Rollback()
		return
	}
	question.AnswerContent = content
	err = databases.UpdateQuestionWithTransaction(ctx, tx, question)
	if err != nil {
		log.Println(err)
		httpError = common.NewInternalServerError(err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func DeleteQuestion(ctx context.Context, questionId uint) (httpError common.HttpError) {
	err := databases.DeleteQuestion(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpError = common.NewNotFoundError("question not found")
			return
		}
		httpError = common.NewInternalServerError(err.Error())
		return
	}
	return
}
