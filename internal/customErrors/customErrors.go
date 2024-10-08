package customErrors

import (
	"fmt"
	"net/http"
	"runtime"
)


type CodedError struct {
	Message   string
	StatusCode int
}

func (e *CodedError) Error() string {
	return fmt.Sprintf("Error: %s, StatusCode: %d", e.Message, e.StatusCode)
}

func GetFunctionName() string {
	pc, _, _, _ := runtime.Caller(1) // Get the program counter of the caller
	function := runtime.FuncForPC(pc)

	return function.Name()
}

func ErrorMarshal(w *http.ResponseWriter, errMarshal error) {
	(*w).WriteHeader(http.StatusInternalServerError)
	fmt.Printf("Error marshalling JSON: %v", errMarshal)
	(*w).WriteHeader(http.StatusInternalServerError)
	(*w).Write(nil)
}