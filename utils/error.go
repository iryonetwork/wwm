package utils

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime"

	authModels "github.com/iryonetwork/wwm/gen/auth/models"
	waitlistModels "github.com/iryonetwork/wwm/gen/waitlist/models"
)

// error codes
const (
	ErrNotFound    = "not_found"
	ErrServerError = "server_error"
	ErrBadRequest  = "bad_request"
	ErrForbidden   = "forbidden"
	ErrConflict    = "conflict"
)

// Error wraps models.Error so it will implement error interface
type Error struct {
	e interface{}
}

// Error returns error message
func (err Error) Error() string {
	switch e := err.e.(type) {
	case authModels.Error:
		return e.Message
	case waitlistModels.Error:
		return e.Message
	}
	return ""
}

// Code returns errors code
func (err Error) Code() string {
	switch e := err.e.(type) {
	case authModels.Error:
		return e.Code
	case waitlistModels.Error:
		return e.Code
	}
	return ""
}

// WriteResponse uses producer to write the error as http response
func (err Error) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	switch err.Code() {
	case ErrBadRequest:
		rw.WriteHeader(400)
	case ErrForbidden:
		rw.WriteHeader(403)
	case ErrNotFound:
		rw.WriteHeader(404)
	case ErrConflict:
		rw.WriteHeader(409)
	default:
		rw.WriteHeader(500)
		err.e = authModels.Error{
			Code:    ErrServerError,
			Message: "Internal Server Error",
		}
	}

	if err.e != nil {
		payload := err.e
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// NewError returns new Error object
func NewError(code, message string, a ...interface{}) Error {
	return Error{
		e: authModels.Error{
			Code:    code,
			Message: fmt.Sprintf(message, a...),
		},
	}
}

// NewServerError creates new server error
func NewServerError(err error) Error {
	return NewError(ErrServerError, err.Error())
}

// NewErrorResponse takes error or Error and returns Error
func NewErrorResponse(e interface{}) Error {
	if err, ok := e.(Error); ok {
		return err
	}

	return NewServerError(e.(error))
}
