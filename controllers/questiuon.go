package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"blog/common"
	"blog/middleware"
	"blog/models"
	"blog/services"
)

func init() {
	Router.Use(Cors())
	questionGroup := Router.Group("/question")
	{
		questionGroup.GET("", getQuestions)
		questionGroup.POST("", newQuestion)
		questionGroup.POST("/:questionId/answer", middleware.AuthMiddleware.MiddlewareFunc(), newAnswer)
	}
}

func getQuestions(c *gin.Context) {
	pager := common.NewPager()
	err := c.BindQuery(pager)
	if err != nil {
		handleError(c, common.NewBadRequestError("error binding pager"))
	}
	questionList, httpError := services.GetQuestions(c, pager.Page, pager.Size)
	if err != nil {
		handleError(c, httpError)
	}
	c.JSON(http.StatusOK, questionList)
}

func newQuestion(c *gin.Context) {
	question := &models.Question{}
	if err := c.ShouldBindJSON(question); err != nil {
		handleError(c, common.NewBadRequestError("error binding request body"))
	}
	if httpError := services.AddQuestion(c, question); httpError != nil {
		handleError(c, httpError)
	}
}

type answerRequest struct {
	answerContent string
}

func newAnswer(c *gin.Context) {
	answer := &answerRequest{}
	if err := c.Bind(answer); err != nil {
		handleError(c, common.NewBadRequestError("error parsing answerContent"))
	}
	questionId, err := strconv.Atoi(c.Param("questionId"))
	if err != nil {
		handleError(c, common.NewBadRequestError("error parsing questionId"))
	}
	httpError := services.AnswerQuestion(c, uint(questionId), answer.answerContent)
	if httpError != nil {
		handleError(c, httpError)
	}
}

func deleteQuestion(c *gin.Context) {

}

func editAnswer() {

}
