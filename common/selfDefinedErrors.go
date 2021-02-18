package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError interface {
	error
	Status() int
	Message() string
}

type SelfDefinedError struct {
	status  int
	message string
}

func NewSelfDefinedError(status int, message string) *SelfDefinedError {
	return &SelfDefinedError{
		status:  status,
		message: message,
	}
}

func (err *SelfDefinedError) Error() string {
	return fmt.Sprintf("status:%d\tmessage:%s", err.Status(), err.Message())
}

func (err *SelfDefinedError) Status() int {
	return err.status
}

func (err *SelfDefinedError) Message() string {
	return err.message
}

func NewNotFoundError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusNotFound, message)
}

func NewInternalServerError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusInternalServerError, message)
}

func NewBadRequestError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusBadRequest, message)
}

func NewInternalError(ctx *gin.Context, err error) {
	_ = ctx.Error(NewInternalServerError("something bad happens")).SetType(gin.ErrorTypePublic)
	_ = ctx.Error(err).SetType(gin.ErrorTypePrivate)
}
