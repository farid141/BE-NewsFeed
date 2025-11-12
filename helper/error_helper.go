package helper

type ServiceError struct {
	StatusCode int
	Message    string
	Details    any
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(status int, msg string, details any) *ServiceError {
	return &ServiceError{StatusCode: status, Message: msg, Details: details}
}
