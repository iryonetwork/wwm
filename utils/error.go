package utils

import (
	"fmt"

	"github.com/iryonetwork/wwm/gen/models"
)

// error codes
const (
	ErrNotFound    = "not_found"
	ErrServerError = "server_error"
	ErrBadRequest  = "bad_request"
)

// Error wraps models.Error so it will implement error interface
type Error struct {
	e models.Error
}

func (err Error) Error() string {
	return err.e.Message
}

// Code returns errors code
func (err Error) Code() string {
	return err.e.Code
}

// NewError returns new Error object
func NewError(code, message string, a ...interface{}) Error {
	return Error{
		e: models.Error{
			Code:    code,
			Message: fmt.Sprintf(message, a...),
		},
	}
}
