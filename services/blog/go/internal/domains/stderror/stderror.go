package stderror

import (
	"errors"
	"net/http"
)

type Code int

func (t Code) HTTPStatusCode() int {
	switch t {
	case CodeBadRequest:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

const (
	CodeUnknown Code = iota
	CodeBadRequest
)

type stdError struct {
	code    Code
	message string
}

func (t stdError) Error() string {
	return t.message
}

func ToCode(err error) Code {
	var stderr *stdError
	if errors.As(err, &stderr) {
		return stderr.code
	}
	return CodeUnknown
}

func NewBadRequest(msg string) *stdError {
	return &stdError{
		code:    CodeBadRequest,
		message: msg,
	}
}
