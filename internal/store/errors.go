package store

import "net/http"

// AppError é um erro de negócio associado a um status HTTP específico,
// permitindo que a camada de handlers traduza regras violadas em respostas.
type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

func ErrNotFound(message string) *AppError {
	return NewAppError(http.StatusNotFound, message)
}

func ErrConflict(message string) *AppError {
	return NewAppError(http.StatusConflict, message)
}

func ErrUnprocessable(message string) *AppError {
	return NewAppError(http.StatusUnprocessableEntity, message)
}

func ErrBadRequest(message string) *AppError {
	return NewAppError(http.StatusBadRequest, message)
}
