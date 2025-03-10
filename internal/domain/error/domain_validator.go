package error

import "strings"

type DomainError struct {
	Type             string
	Operation        string
	Message          string
	Err              error
	ValidationErrors []error
}

// NewDomainError crea un nuevo error de dominio con el tipo de servicio, la operación, el mensaje y el error.
func NewDomainError(serviceType, op, msg string) *DomainError {
	return &DomainError{
		Type:      serviceType,
		Operation: op,
		Message:   msg,
	}
}

func NewDomainErrorWithCause(serviceType, op, msg string, err error) *DomainError {
	return &DomainError{
		Type:      serviceType,
		Operation: op,
		Message:   msg,
		Err:       err,
	}
}

func (e *DomainError) AsideError(err error) {
	e.Err = err
}

// Error Implementación de la interfaz error para el error de dominio para lanzar error principal a nivel de servicio y errores de validación
func (e *DomainError) Error() string {
	if e.Err == nil {
		return e.Message
	}

	return e.Message + " | cause: " + e.Err.Error()
}

func (e *DomainError) IsNotFoundError() bool {
	return strings.Contains(e.Error(), "not found")
}

func (e *DomainError) HasValidationErrors() bool {
	return len(e.ValidationErrors) > 0
}

// AddValidationError agrega un error de validación al error de dominio
func (e *DomainError) AddValidationError(err error) {
	e.ValidationErrors = append(e.ValidationErrors, err)
}

// GetValidationErrorsString Obtiene los errores de validación asociados al error de dominio en caso de existir
func (e *DomainError) GetValidationErrorsString() []string {
	var messages []string
	for _, err := range e.ValidationErrors {
		messages = append(messages, err.Error())
	}

	return messages
}
