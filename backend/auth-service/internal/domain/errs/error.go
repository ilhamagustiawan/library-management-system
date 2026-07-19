package errs

import "errors"

const (
	CodeSuccess            = "LMS-200000"
	CodeValidation         = "LMS-422001"
	CodeInvalidCredentials = "LMS-401001"
	CodeInvalidToken       = "LMS-401002"
	CodeUserNotFound       = "LMS-404001"
	CodeEmailExists        = "LMS-409001"
	CodeRateLimited        = "LMS-429001"
	CodeNotFound           = "LMS-404000"
	CodeInternal           = "LMS-500000"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record conflict")
)

type Error struct {
	HTTPStatus int
	ErrorCode  string
	Message    string
	Data       any
	Cause      error
}

func (e *Error) Error() string { return e.Message }

func (e *Error) Unwrap() error { return e.Cause }

func New(httpStatus int, code, message string, data any, cause error) *Error {
	return &Error{
		HTTPStatus: httpStatus,
		ErrorCode:  code,
		Message:    message,
		Data:       data,
		Cause:      cause,
	}
}
