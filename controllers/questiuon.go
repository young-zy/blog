package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/young-zy/blog/common"
	"github.com/young-zy/blog/middleware"
	"github.com/young-zy/blog/models"
	"github.com/young-zy/blog/services"
)

func initQuestionGroup() {
	questionGroup := Router.Group("/question")
	{
		questionGroup.GET("", getQuestions)
		questionGroup.GET("/:questionID", getQuestion)
		questionGroup.POST("", middleware.Recaptcha.Middleware(), newQuestion)
		questionGroup.POST("/:questionID/answer", middleware.AuthMiddleware.MiddlewareFunc(), newAnswer)
		questionGroup.DELETE("/:questionID", middleware.AuthMiddleware.MiddlewareFunc(), deleteQuestion)
		questionGroup.PATCH("/:questionID", middleware.AuthMiddleware.MiddlewareFunc(), updateAnswer)
	}
}

func getQuestion(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionID")).SetType(gin.ErrorTypeBind)
	}
	if question, ok := services.GetQuestion(c, common.IntToUintPointer(questionID)); ok {
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
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionID")).SetType(gin.ErrorTypeBind)
		return
	}
	if services.AnswerQuestion(c, common.IntToUintPointer(questionID), &answer.AnswerContent) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func deleteQuestion(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionID")).SetType(gin.ErrorTypeBind)
		return
	}
	if services.DeleteQuestion(c, common.IntToUintPointer(questionID)) {
		c.JSON(http.StatusNoContent, nil)
	}
}

func updateAnswer(c *gin.Context) {
	questionID, err := strconv.Atoi(c.Param("questionID"))
	if err != nil {
		_ = c.Error(errors.New("error parsing questionID")).SetType(gin.ErrorTypeBind)
		return
	}
	question := &models.Question{}
	if c.Bind(question) != nil {
		return
	}
	services.UpdateAnswer(c, common.IntToUintPointer(questionID), question.AnswerContent)
	if services.DeleteQuestion(c, common.IntToUintPointer(questionID)) {
		c.JSON(http.StatusNoContent, nil)
	}
}
