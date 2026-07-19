package errs

import "errors"

const (
	CodeSuccess             = "LMS-200000"
	CodeUnauthorized        = "LMS-401001"
	CodeForbidden           = "LMS-403001"
	CodeNotFound            = "LMS-404000"
	CodeBookNotFound        = "LMS-404002"
	CodeReservationMissing  = "LMS-404003"
	CodeISBNExists          = "LMS-409002"
	CodeStockUnavailable    = "LMS-409003"
	CodeInventoryConflict   = "LMS-409004"
	CodeReservationConflict = "LMS-409005"
	CodeValidation          = "LMS-422001"
	CodeAuthUnavailable     = "LMS-503001"
	CodeInternal            = "LMS-500000"
)

var (
	ErrNotFound         = errors.New("record not found")
	ErrConflict         = errors.New("record conflict")
	ErrISBNExists       = errors.New("ISBN already exists")
	ErrStockUnavailable = errors.New("book stock unavailable")
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
