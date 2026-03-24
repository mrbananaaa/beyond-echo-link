package apperror

const (
	CodeInternal     = "INTERNAL_ERROR"
	CodeBadRequest   = "BAD_REQUEST"
	CodeNotFound     = "NOT_FOUND"
	CodeUnauthorized = "UNAUTHORIZED"
)

var (
	ErrInternal     = New(CodeInternal, "internal server error")
	ErrBadRequest   = New(CodeBadRequest, "bad request")
	ErrNotFound     = New(CodeNotFound, "resource not found")
	ErrUnauthorized = New(CodeUnauthorized, "unauthorized")
)

type AppEror struct {
	Code    string
	Message string
	Details []string
}

func (e *AppEror) Error() string {
	return e.Message
}

func New(code, message string, details ...string) *AppEror {
	return &AppEror{
		Code:    code,
		Message: message,
		Details: details,
	}
}
