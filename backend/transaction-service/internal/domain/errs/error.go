package errs

import "errors"

const (
	CodeSuccess          = "LMS-200000"
	CodeValidation       = "LMS-422001"
	CodeInvalidToken     = "LMS-401002"
	CodeForbidden        = "LMS-403001"
	CodeLoanNotFound     = "LMS-404004"
	CodeBookNotFound     = "LMS-404002"
	CodeLoanLimit        = "LMS-409004"
	CodeActiveLoan       = "LMS-409005"
	CodeStockUnavailable = "LMS-409003"
	CodeDependency       = "LMS-503001"
	CodeInternal         = "LMS-500000"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrForbidden        = errors.New("operation forbidden")
	ErrLoanLimit        = errors.New("active loan limit reached")
	ErrActiveLoan       = errors.New("book already has an active loan")
	ErrBookNotFound     = errors.New("book not found")
	ErrStockUnavailable = errors.New("book stock unavailable")
	ErrDependency       = errors.New("dependency unavailable")
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

func New(status int, code, message string, data any, cause error) *Error {
	return &Error{HTTPStatus: status, ErrorCode: code, Message: message, Data: data, Cause: cause}
}
