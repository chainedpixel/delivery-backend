package dto

import (
	"encoding/json"
	"fmt"
	"io"
)

// EmailTestRequest representa la solicitud para enviar un email de prueba
type EmailTestRequest struct {
	// Email de destino
	To string `json:"to" example:"test@example.com" binding:"required,email"`

	// Asunto del email
	Subject string `json:"subject" example:"Test Email Subject" binding:"required"`

	// Cuerpo del email (puede ser HTML o texto plano)
	Body string `json:"body" example:"<h1>This is a test email</h1><p>Hello from delivery system!</p>" binding:"required"`

	// Tipo de contenido (html o text)
	ContentType string `json:"content_type" example:"html" binding:"required"`

	// CC (opcional)
	CC []string `json:"cc,omitempty" example:"[\"cc1@example.com\", \"cc2@example.com\"]"`

	// BCC (opcional)
	BCC []string `json:"bcc,omitempty" example:"[\"bcc@example.com\"]"`

	// Adjuntos (opcional) - URLs o paths a archivos
	Attachments []string `json:"attachments,omitempty" example:"[\"path/to/file.pdf\"]"`
}

// EmailTestResponse representa la respuesta del envío de email de prueba
type EmailTestResponse struct {
	// Estado del envío
	Success bool `json:"success" example:"true"`

	// Mensaje de resultado
	Message string `json:"message" example:"Email sent successfully"`

	// ID del mensaje (si está disponible)
	MessageID string `json:"message_id,omitempty" example:"msg_12345"`

	// Detalles adicionales
	Details map[string]interface{} `json:"details,omitempty"`
}

// NewEmailTestRequest crea una nueva instancia desde el body de la request
func NewEmailTestRequest(body io.ReadCloser) (*EmailTestRequest, error) {
	var request EmailTestRequest
	err := json.NewDecoder(body).Decode(&request)
	if err != nil {
		return nil, fmt.Errorf("failed to decode email test request: %w", err)
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}

// Validate valida los campos requeridos
func (r *EmailTestRequest) Validate() error {
	if r.To == "" {
		return fmt.Errorf("recipient email is required")
	}

	if r.Subject == "" {
		return fmt.Errorf("email subject is required")
	}

	if r.Body == "" {
		return fmt.Errorf("email body is required")
	}

	if r.ContentType == "" {
		r.ContentType = "html"
	}

	// Validar que content_type sea válido
	if r.ContentType != "html" && r.ContentType != "text" {
		return fmt.Errorf("content_type must be 'html' or 'text'")
	}

	return nil
}
