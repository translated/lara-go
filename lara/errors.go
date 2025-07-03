package lara

import "fmt"

type LaraError struct {
	Status  int
	Type    string
	Message string
}

func (e *LaraError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type LaraConnectionError struct {
	Message string
}

func (e *LaraConnectionError) Error() string {
	return fmt.Sprintf("ConnectionError: %s", e.Message)
}

type LaraTimeoutError struct {
	Message string
}

func (e *LaraTimeoutError) Error() string {
	return fmt.Sprintf("TimeoutError: %s", e.Message)
}
