package helper

const (
	CodeValidationError = "VALIDATION_ERROR"
	CodeNotFound        = "NOT_FOUND"
	CodeDuplicate       = "DUPLICATE"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
	CodeConflict        = "CONFLICT"
	CodeInternal        = "INTERNAL_ERROR"
	CodeBadRequest      = "BAD_REQUEST"
	CodeNotAcceptable   = "NOT_ACCEPTABLE"
)

type AppError struct {
	Code    string
	Message string
	Status  int
	Details interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code string, message string, status int, details interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}

func ErrValidation(details interface{}) *AppError {
	return NewAppError(CodeValidationError, "input tidak valid", 400, details)
}

func ErrNotFound(message string) *AppError {
	return NewAppError(CodeNotFound, message, 404, nil)
}

func ErrDuplicate(message string) *AppError {
	return NewAppError(CodeDuplicate, message, 409, nil)
}

func ErrUnauthorized(message string) *AppError {
	return NewAppError(CodeUnauthorized, message, 401, nil)
}

func ErrForbidden(message string) *AppError {
	return NewAppError(CodeForbidden, message, 403, nil)
}

func ErrConflict(message string) *AppError {
	return NewAppError(CodeConflict, message, 409, nil)
}

func ErrInternal(message string) *AppError {
	return NewAppError(CodeInternal, message, 500, nil)
}
