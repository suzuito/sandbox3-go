package repoerror

import "errors"

type Code int

const (
	CodeUnknown Code = iota
	CodeNotFound
)

type repoError struct {
	code    Code
	message string
}

func (t *repoError) Error() string {
	return t.message
}

func ToCode(err error) Code {
	var stderr *repoError
	if errors.As(err, &stderr) {
		return stderr.code
	}
	return CodeUnknown
}

func Has(err error, code Code) bool {
	return ToCode(err) == code
}

func NewNotFound(msg string) *repoError {
	return &repoError{
		code:    CodeNotFound,
		message: msg,
	}
}
