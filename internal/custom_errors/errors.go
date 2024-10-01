package custom_errors

import "fmt"

type CredentialsError struct {
	Message string
}

func (e *CredentialsError) Error() string {
	return e.Message
}

func NewCredentialsError(message string, args ...any) *CredentialsError {
	return &CredentialsError{
		Message: fmt.Sprintf(message, args...),
	}
}
