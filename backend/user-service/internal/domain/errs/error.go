package errs

import "errors"

const (
	CodeSuccess             = "LMS-200000"
	CodeValidation          = "LMS-422001"
	CodeEmailExists         = "LMS-409001"
	CodeIdempotencyConflict = "LMS-409002"
	CodeRateLimited         = "LMS-429001"
	CodeDependency          = "LMS-503001"
	CodeNotFound            = "LMS-404000"
	CodeInternal            = "LMS-500000"
)

var (
	ErrConflict            = errors.New("record conflict")
	ErrIdentityConflict    = errors.New("identity conflict")
	ErrIdentityIdempotency = errors.New("identity idempotency conflict")
	ErrIdentityUnavailable = errors.New("identity service unavailable")
	ErrInvalidIdentity     = errors.New("invalid identity response")
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
	return &Error{HTTPStatus: httpStatus, ErrorCode: code, Message: message, Data: data, Cause: cause}
}
