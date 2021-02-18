package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog/common"
	"blog/middleware"
	"blog/models"
	"blog/services"
)

func initQuestionGroup() {
	questionGroup := Router.Group("/question")
	{
		questionGroup.GET("", getQuestions)
		questionGroup.POST("", newQuestion)
		questionGroup.POST("/:questionId/answer", middleware.AuthMiddleware.MiddlewareFunc(), newAnswer)
		questionGroup.DELETE("/:questionId", middleware.AuthMiddleware.MiddlewareFunc(), deleteQuestion)
		questionGroup.PATCH("/:questionId", middleware.AuthMiddleware.MiddlewareFunc(), updateAnswer)
	}
}

func getQuestions(c *gin.Context) {
	pager := common.NewPager()
	if c.BindQuery(pager) != nil {
		return
	}
	if questionList, ok := services.GetQuestions(c, pager.Page, pager.Size); ok {
		c.JSON(http.StatusOK, questionList)
	}
}

func newQuestion(c *gin.Context) {
	question := &models.Question{}
	if c.ShouldBindJSON(question) != nil {
		return
	}
	services.AddQuestion(c, question)
}

type answerRequest struct {
	answerContent string
}

func newAnswer(c *gin.Context) {
	answer := &answerRequest{}
	if c.Bind(answer) != nil {
		return
	}
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionId")).SetType(gin.ErrorTypeBind)
		return
	}
	if services.AnswerQuestion(c, uint(questionId), answer.answerContent) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func deleteQuestion(c *gin.Context) {
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionId")).SetType(gin.ErrorTypeBind)
		return
	}
	if services.DeleteQuestion(c, uint(questionId)) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func updateAnswer(c *gin.Context) {
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionId")).SetType(gin.ErrorTypeBind)
		return
	}
	question := &models.Question{}
	if c.Bind(question) != nil {
		return
	}
	services.UpdateAnswer(c, uint(questionId), question.AnswerContent)
	if services.DeleteQuestion(c, uint(questionId)) {
		c.JSON(http.StatusNoContent, nil)
	}
}
