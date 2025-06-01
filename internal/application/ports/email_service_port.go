package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

// EmailData contiene la información necesaria para enviar un email
type EmailData struct {
	To          []string               // Destinatarios
	CC          []string               // Copia
	BCC         []string               // Copia oculta
	Subject     string                 // Asunto
	Body        string                 // Cuerpo del email (HTML)
	PlainText   string                 // Cuerpo alternativo en texto plano
	Attachments []EmailAttachment      // Archivos adjuntos
	Variables   map[string]interface{} // Variables para templates
}

// EmailAttachment representa un archivo adjunto
type EmailAttachment struct {
	Filename string // Nombre del archivo
	Content  []byte // Contenido del archivo
	MimeType string // Tipo MIME
}

// EmailTemplate contiene la información de un template de email
type EmailTemplate struct {
	Name      string   // Nombre del template
	Subject   string   // Asunto del email
	HTMLBody  string   // Cuerpo HTML
	PlainBody string   // Cuerpo texto plano
	Variables []string // Variables disponibles
}

// OrderEmailData contiene los datos específicos para emails de pedidos
type OrderEmailData struct {
	Order           *entities.Order
	Customer        *entities.User
	Company         *entities.Company
	TrackingURL     string
	EstimatedTime   string
	DriverInfo      *DriverEmailInfo
	AdditionalNotes string
}

// DriverEmailInfo contiene información del conductor para el email
type DriverEmailInfo struct {
	Name        string
	Phone       string
	VehicleInfo string
	Photo       string
}

// EmailService define las operaciones para el envío de correos electrónicos
type EmailService interface {
	// SendEmail envía un email genérico
	SendEmail(ctx context.Context, emailData *EmailData) error
	// SendTemplatedEmail envía un email usando un template
	SendTemplatedEmail(ctx context.Context, templateName string, to []string, variables map[string]interface{}) error
	// SendOrderEmail envía un email relacionado con un pedido
	SendOrderEmail(ctx context.Context, emailType string, orderData *OrderEmailData) error
	// ValidateEmailTemplate valida que un template sea correcto
	ValidateEmailTemplate(template *EmailTemplate) error
	// GetEmailTemplate obtiene un template por nombre
	GetEmailTemplate(templateName string) (*EmailTemplate, error)
	// SetEmailTemplate establece o actualiza un template
	SetEmailTemplate(template *EmailTemplate) error
	SendTestEmail(ctx context.Context, to, subject, body, contentType string, cc, bcc []string) (messageID string, err error)
	GetServiceStatus() map[string]interface{}
	SendWelcomeTestEmail(ctx context.Context, to, name string) (messageID string, err error)
}
