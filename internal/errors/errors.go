package errors

import (
	"fmt"
)


type CodedError struct {
	Message   string
	StatusCode int
}

func (e *CodedError) Error() string {
	return fmt.Sprintf("Error: %s, StatusCode: %d", e.Message, e.StatusCode)
}