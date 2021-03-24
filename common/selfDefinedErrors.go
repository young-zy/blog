package common

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HttpError is an interface used to present an http error
type HttpError interface {
	error
	Status() int
	Message() string
}

// SelfDefinedError is an implementation of HttpError
type SelfDefinedError struct {
	status  int
	message string
}

// NewSelfDefinedError creates a SelfDefinedError
func NewSelfDefinedError(status int, message string) *SelfDefinedError {
	return &SelfDefinedError{
		status:  status,
		message: message,
	}
}

// Error returns error
func (err *SelfDefinedError) Error() string {
	return fmt.Sprintf("status:%d\tmessage:%s", err.Status(), err.Message())
}

// Status returns status
func (err *SelfDefinedError) Status() int {
	return err.status
}

// Message returns message
func (err *SelfDefinedError) Message() string {
	return err.message
}

// NewNotFoundError is a shortcut to create a not found http error
func NewNotFoundError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusNotFound, message)
}

// NewInternalServerError is a shortcut to create an internal http server error
func NewInternalServerError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusInternalServerError, message)
}

// NewBadRequestError is a shortcut to create a bad request http error
func NewBadRequestError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusBadRequest, message)
}

// NewForbiddenError is a shortcut to create a forbidden http error
func NewForbiddenError(message string) *SelfDefinedError {
	return NewSelfDefinedError(http.StatusForbidden, message)
}

// NewInternalError keeps the private error to gin and creates an internal server error
func NewInternalError(ctx *gin.Context, err error) {
	_ = ctx.Error(NewInternalServerError("something bad happens")).SetType(gin.ErrorTypePublic)
	_ = ctx.Error(err).SetType(gin.ErrorTypePrivate)
}
