package error

import (
	"config"
	"fmt"
)

type ServiceError struct {
	Type      string
	Operation string
	Err       error
}

func (e *ServiceError) Error() string {
	envConfig, err := config.NewEnvConfig()
	if err != nil {
		return ""
	}
	if envConfig.Log.Level == "debug" {
		return fmt.Sprintf("[%s] %s: | cause: %v", e.Type, e.Operation, e.Err.Error())
	}

	return fmt.Sprintf("%s", e.Err.Error())
}

// NewGeneralServiceError crea un nuevo error de servicio general con el tipo de servicio, la operación, el mensaje y el error.
func NewGeneralServiceError(serviceType, op string, err error) *ServiceError {
	return &ServiceError{
		Type:      serviceType,
		Operation: op,
		Err:       err,
	}
}
