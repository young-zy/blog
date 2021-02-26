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
		questionGroup.GET("/:questionId", getQuestion)
		questionGroup.POST("", middleware.Recaptcha.Middleware(), newQuestion)
		questionGroup.POST("/:questionId/answer", middleware.AuthMiddleware.MiddlewareFunc(), newAnswer)
		questionGroup.DELETE("/:questionId", middleware.AuthMiddleware.MiddlewareFunc(), deleteQuestion)
		questionGroup.PATCH("/:questionId", middleware.AuthMiddleware.MiddlewareFunc(), updateAnswer)
	}
}

func getQuestion(c *gin.Context) {
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionId")).SetType(gin.ErrorTypeBind)
	}
	if question, ok := services.GetQuestion(c, common.IntToUintPointer(questionId)); ok {
		c.JSON(http.StatusOK, question)
	}
}

// default page is 1 and default size is 10.
func getQuestions(c *gin.Context) {
	pager := common.NewPager()
	if c.BindQuery(pager) != nil {
		return
	}
	filter := c.Query("filter")
	if questionList, ok := services.GetQuestions(c, pager.Page, pager.Size, filter); ok {
		c.JSON(http.StatusOK, questionList)
	}
}

func newQuestion(c *gin.Context) {
	question := &models.NewQuestionRequest{}
	if c.BindJSON(question) != nil {
		return
	}
	services.AddQuestion(c, question)
}

type answerRequest struct {
	AnswerContent string `binding:"required" json:"answerContent"`
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
	if services.AnswerQuestion(c, common.IntToUintPointer(questionId), &answer.AnswerContent) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func deleteQuestion(c *gin.Context) {
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionId")).SetType(gin.ErrorTypeBind)
		return
	}
	if services.DeleteQuestion(c, common.IntToUintPointer(questionId)) {
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
	services.UpdateAnswer(c, common.IntToUintPointer(questionId), question.AnswerContent)
	if services.DeleteQuestion(c, common.IntToUintPointer(questionId)) {
		c.JSON(http.StatusNoContent, nil)
	}
}
