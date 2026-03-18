package apperror

const (
	CodeInternal   = "INTERNAL_ERROR"
	CodeBadRequest = "BAD_REQUEST"
	CodeNotFound   = "NOT_FOUND"
)

type Error struct {
	Code    string
	Message string
	Details []string
}

func (e *Error) Error() string {
	return e.Message
}

func New(code, message string, details ...string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func ErrNotFound(msg string) *Error {
	return New(CodeNotFound, msg)
}

func ErrInternal() *Error {
	return New(CodeInternal, "internal server error")
}

func ErrBadRequest(msg string, details ...string) *Error {
	return New(CodeBadRequest, msg, details...)
}
