package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"strings"
	"time"

	"github.com/MarlonG1/delivery-backend/configs"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	errPackage "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	_ "gopkg.in/gomail.v2"
)

type EmailServiceImpl struct {
	config    *config.EnvConfig
	templates map[string]*ports.EmailTemplate
	dialer    *gomail.Dialer
}

// NewEmailService crea una nueva instancia del servicio de email
func NewEmailService(emailConfig *config.EnvConfig) ports.EmailService {
	// Configurar el dialer SMTP
	dialer := gomail.NewDialer(
		emailConfig.EmailConfig.SMTPHost,
		emailConfig.EmailConfig.SMTPPort,
		emailConfig.EmailConfig.SMTPUsername,
		emailConfig.EmailConfig.SMTPPassword,
	)

	// Configurar TLS si está habilitado
	if emailConfig.EmailConfig.EnableTLS {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: emailConfig.EmailConfig.SkipTLSVerify,
			ServerName:         emailConfig.EmailConfig.SMTPHost,
		}
	}

	service := &EmailServiceImpl{
		config:    emailConfig,
		templates: make(map[string]*ports.EmailTemplate),
		dialer:    dialer,
	}

	// Cargar templates por defecto
	service.loadDefaultTemplates()

	return service
}

// SendEmail envía un email genérico
func (s *EmailServiceImpl) SendEmail(ctx context.Context, emailData *ports.EmailData) error {
	if len(emailData.To) == 0 {
		return errPackage.NewGeneralServiceError("EmailService", "SendEmail", errPackage.ErrNoRecipients)
	}

	message := gomail.NewMessage()

	// Configurar remitente
	message.SetHeader("From", s.config.EmailConfig.FromEmail, s.config.EmailConfig.FromEmail)

	// Configurar destinatarios
	message.SetHeader("To", emailData.To...)
	if len(emailData.CC) > 0 {
		message.SetHeader("Cc", emailData.CC...)
	}
	if len(emailData.BCC) > 0 {
		message.SetHeader("Bcc", emailData.BCC...)
	}

	// Configurar asunto
	message.SetHeader("Subject", emailData.Subject)

	// Configurar cuerpo
	if emailData.Body != "" {
		if emailData.PlainText != "" {
			message.SetBody("text/plain", emailData.PlainText)
			message.AddAlternative("text/html", emailData.Body)
		} else {
			message.SetBody("text/html", emailData.Body)
		}
	} else if emailData.PlainText != "" {
		message.SetBody("text/plain", emailData.PlainText)
	}

	// Agregar archivos adjuntos
	for _, attachment := range emailData.Attachments {
		message.Attach(attachment.Filename,
			gomail.SetHeader(map[string][]string{
				"Content-Type": {attachment.MimeType},
			}))
	}

	// Enviar email
	if err := s.dialer.DialAndSend(message); err != nil {
		logs.Error("Failed to send email", map[string]interface{}{
			"error":   err.Error(),
			"to":      emailData.To,
			"subject": emailData.Subject,
		})
		return errPackage.NewGeneralServiceError("EmailService", "SendEmail", err)
	}

	logs.Info("Email sent successfully", map[string]interface{}{
		"to":      emailData.To,
		"subject": emailData.Subject,
	})

	return nil
}

// SendTemplatedEmail envía un email usando un template
func (s *EmailServiceImpl) SendTemplatedEmail(ctx context.Context, templateName string, to []string, variables map[string]interface{}) error {
	emailTemplate, exists := s.templates[templateName]
	if !exists {
		return errPackage.NewGeneralServiceError("EmailService", "SendTemplatedEmail",
			fmt.Errorf("template %s not found", templateName))
	}

	// Procesar template HTML
	htmlBody, err := s.processTemplate(emailTemplate.HTMLBody, variables)
	if err != nil {
		return errPackage.NewGeneralServiceError("EmailService", "SendTemplatedEmail", err)
	}

	// Procesar template de texto plano
	plainBody, err := s.processTemplate(emailTemplate.PlainBody, variables)
	if err != nil {
		return errPackage.NewGeneralServiceError("EmailService", "SendTemplatedEmail", err)
	}

	// Procesar asunto
	subject, err := s.processTemplate(emailTemplate.Subject, variables)
	if err != nil {
		return errPackage.NewGeneralServiceError("EmailService", "SendTemplatedEmail", err)
	}

	emailData := &ports.EmailData{
		To:        to,
		Subject:   subject,
		Body:      htmlBody,
		PlainText: plainBody,
		Variables: variables,
	}

	return s.SendEmail(ctx, emailData)
}

// SendOrderEmail envía un email relacionado con un pedido
func (s *EmailServiceImpl) SendOrderEmail(ctx context.Context, emailType string, orderData *ports.OrderEmailData) error {
	templateName := s.getOrderEmailTemplate(emailType)

	variables := s.buildOrderEmailVariables(orderData)

	// Obtener emails de los destinatarios
	recipients := []string{orderData.Customer.Email}

	// Agregar email de la empresa si es diferente
	if orderData.Company != nil && orderData.Company.ContactEmail != orderData.Customer.Email {
		recipients = append(recipients, orderData.Company.ContactEmail)
	}

	return s.SendTemplatedEmail(ctx, templateName, recipients, variables)
}

// ValidateEmailTemplate valida que un template sea correcto
func (s *EmailServiceImpl) ValidateEmailTemplate(emailTemplate *ports.EmailTemplate) error {
	if emailTemplate.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if emailTemplate.Subject == "" {
		return fmt.Errorf("template subject is required")
	}

	if emailTemplate.HTMLBody == "" && emailTemplate.PlainBody == "" {
		return fmt.Errorf("template must have either HTML or plain text body")
	}

	// Validar que los templates sean válidos
	_, err := template.New("test").Parse(emailTemplate.HTMLBody)
	if err != nil {
		return fmt.Errorf("invalid HTML template: %w", err)
	}

	if emailTemplate.PlainBody != "" {
		_, err := template.New("test").Parse(emailTemplate.PlainBody)
		if err != nil {
			return fmt.Errorf("invalid plain text template: %w", err)
		}
	}

	return nil
}

// GetEmailTemplate obtiene un template por nombre
func (s *EmailServiceImpl) GetEmailTemplate(templateName string) (*ports.EmailTemplate, error) {
	emailTemplate, exists := s.templates[templateName]
	if !exists {
		return nil, fmt.Errorf("template %s not found", templateName)
	}
	return emailTemplate, nil
}

// SetEmailTemplate establece o actualiza un template
func (s *EmailServiceImpl) SetEmailTemplate(emailTemplate *ports.EmailTemplate) error {
	if err := s.ValidateEmailTemplate(emailTemplate); err != nil {
		return err
	}

	s.templates[emailTemplate.Name] = emailTemplate
	return nil
}

// processTemplate procesa un template con las variables proporcionadas
func (s *EmailServiceImpl) processTemplate(templateStr string, variables map[string]interface{}) (string, error) {
	if templateStr == "" {
		return "", nil
	}

	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, variables); err != nil {
		return "", err
	}

	return result.String(), nil
}

// getOrderEmailTemplate obtiene el nombre del template según el tipo de email de pedido
func (s *EmailServiceImpl) getOrderEmailTemplate(emailType string) string {
	switch emailType {
	case "order_created":
		return "order_created"
	case "order_started":
		return "order_started"
	case "order_completed":
		return "order_completed"
	case "order_cancelled":
		return "order_cancelled"
	default:
		return "order_update"
	}
}

// buildOrderEmailVariables construye las variables para los templates de pedidos
func (s *EmailServiceImpl) buildOrderEmailVariables(orderData *ports.OrderEmailData) map[string]interface{} {
	variables := map[string]interface{}{
		"CustomerName":    orderData.Customer.FullName,
		"OrderID":         orderData.Order.ID,
		"TrackingNumber":  orderData.Order.TrackingNumber,
		"TrackingURL":     orderData.TrackingURL,
		"OrderStatus":     orderData.Order.Status,
		"CreatedAt":       orderData.Order.CreatedAt.Format("02/01/2006 15:04"),
		"CompanyName":     "",
		"CompanyPhone":    "",
		"CompanyEmail":    "",
		"EstimatedTime":   orderData.EstimatedTime,
		"AdditionalNotes": orderData.AdditionalNotes,
		"CurrentYear":     time.Now().Year(),
	}

	// Información de la empresa
	if orderData.Company != nil {
		variables["CompanyName"] = orderData.Company.Name
		variables["CompanyPhone"] = orderData.Company.ContactPhone
		variables["CompanyEmail"] = orderData.Company.ContactEmail
	}

	// Información del pedido
	if orderData.Order.Detail != nil {
		variables["OrderPrice"] = fmt.Sprintf("%.2f", orderData.Order.Detail.Price)
		variables["PickupTime"] = orderData.Order.Detail.PickupTime.Format("02/01/2006 15:04")
		variables["DeliveryDeadline"] = orderData.Order.Detail.DeliveryDeadline.Format("02/01/2006 15:04")
		variables["DeliveryNotes"] = orderData.Order.Detail.DeliveryNotes
	}

	// Información de dirección de entrega
	if orderData.Order.DeliveryAddress != nil {
		variables["DeliveryAddress"] = fmt.Sprintf("%s, %s, %s",
			orderData.Order.DeliveryAddress.AddressLine1,
			orderData.Order.DeliveryAddress.City,
			orderData.Order.DeliveryAddress.State)
		variables["RecipientName"] = orderData.Order.DeliveryAddress.RecipientName
		variables["RecipientPhone"] = orderData.Order.DeliveryAddress.RecipientPhone
	}

	// Información del conductor
	if orderData.DriverInfo != nil {
		variables["DriverName"] = orderData.DriverInfo.Name
		variables["DriverPhone"] = orderData.DriverInfo.Phone
		variables["VehicleInfo"] = orderData.DriverInfo.VehicleInfo
	}

	return variables
}

// loadDefaultTemplates carga los templates por defecto
func (s *EmailServiceImpl) loadDefaultTemplates() {
	// Template para pedido creado
	s.templates["order_created"] = &ports.EmailTemplate{
		Name:    "order_created",
		Subject: "Pedido Confirmado - {{.TrackingNumber}}",
		HTMLBody: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Pedido Confirmado</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50; border-bottom: 2px solid #3498db; padding-bottom: 10px;">
            ¡Pedido Confirmado!
        </h2>
        
        <p>Hola <strong>{{.CustomerName}}</strong>,</p>
        
        <p>Tu pedido ha sido confirmado y está siendo procesado.</p>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3>Detalles del Pedido:</h3>
            <ul style="list-style: none; padding: 0;">
                <li><strong>Número de Seguimiento:</strong> {{.TrackingNumber}}</li>
                <li><strong>Estado:</strong> {{.OrderStatus}}</li>
                <li><strong>Precio:</strong> ${{.OrderPrice}}</li>
                <li><strong>Fecha de Creación:</strong> {{.CreatedAt}}</li>
            </ul>
        </div>
        
        <div style="background-color: #e8f5e8; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3>Dirección de Entrega:</h3>
            <p><strong>{{.RecipientName}}</strong><br>
            {{.DeliveryAddress}}<br>
            Tel: {{.RecipientPhone}}</p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="{{.TrackingURL}}" 
               style="background-color: #3498db; color: white; padding: 12px 25px; 
                      text-decoration: none; border-radius: 5px; display: inline-block;">
                Rastrear Pedido
            </a>
        </p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        
        <p style="font-size: 12px; color: #666; text-align: center;">
            {{.CompanyName}}<br>
            {{.CompanyPhone}} | {{.CompanyEmail}}<br>
            © {{.CurrentYear}} Todos los derechos reservados
        </p>
    </div>
</body>
</html>`,
		PlainBody: `Hola {{.CustomerName}},

Tu pedido ha sido confirmado y está siendo procesado.

Detalles del Pedido:
- Número de Seguimiento: {{.TrackingNumber}}
- Estado: {{.OrderStatus}}
- Precio: ${{.OrderPrice}}
- Fecha de Creación: {{.CreatedAt}}

Dirección de Entrega:
{{.RecipientName}}
{{.DeliveryAddress}}
Tel: {{.RecipientPhone}}

Puedes rastrear tu pedido en: {{.TrackingURL}}

{{.CompanyName}}
{{.CompanyPhone}} | {{.CompanyEmail}}`,
	}

	// Template para pedido iniciado
	s.templates["order_started"] = &ports.EmailTemplate{
		Name:    "order_started",
		Subject: "Pedido en Camino - {{.TrackingNumber}}",
		HTMLBody: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Pedido en Camino</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #27ae60; border-bottom: 2px solid #27ae60; padding-bottom: 10px;">
¡Tu pedido está en camino!
        </h2>
        
        <p>Hola <strong>{{.CustomerName}}</strong>,</p>
        
        <p>¡Excelentes noticias! Tu pedido ya está en camino y llegará pronto.</p>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3>Información del Pedido:</h3>
            <ul style="list-style: none; padding: 0;">
                <li><strong>Número de Seguimiento:</strong> {{.TrackingNumber}}</li>
                <li><strong>Estado:</strong> {{.OrderStatus}}</li>
                {{if .EstimatedTime}}<li><strong>Tiempo Estimado:</strong> {{.EstimatedTime}}</li>{{end}}
            </ul>
        </div>
        
        {{if .DriverName}}
        <div style="background-color: #e3f2fd; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3>Información del Conductor:</h3>
            <p><strong>{{.DriverName}}</strong><br>
            Tel: {{.DriverPhone}}<br>
            {{if .VehicleInfo}}Vehículo: {{.VehicleInfo}}{{end}}</p>
        </div>
        {{end}}
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="{{.TrackingURL}}" 
               style="background-color: #27ae60; color: white; padding: 12px 25px; 
                      text-decoration: none; border-radius: 5px; display: inline-block;">
                Rastrear en Tiempo Real
            </a>
        </p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        
        <p style="font-size: 12px; color: #666; text-align: center;">
            {{.CompanyName}}<br>
            {{.CompanyPhone}} | {{.CompanyEmail}}<br>
            © {{.CurrentYear}} Todos los derechos reservados
        </p>
    </div>
</body>
</html>`,
		PlainBody: `Hola {{.CustomerName}},

¡Excelentes noticias! Tu pedido ya está en camino y llegará pronto.

Información del Pedido:
- Número de Seguimiento: {{.TrackingNumber}}
- Estado: {{.OrderStatus}}
{{if .EstimatedTime}}- Tiempo Estimado: {{.EstimatedTime}}{{end}}

{{if .DriverName}}Información del Conductor:
{{.DriverName}}
Tel: {{.DriverPhone}}
{{if .VehicleInfo}}Vehículo: {{.VehicleInfo}}{{end}}{{end}}

Puedes rastrear tu pedido en tiempo real en: {{.TrackingURL}}

{{.CompanyName}}
{{.CompanyPhone}} | {{.CompanyEmail}}`,
	}

	// Template para pedido completado
	s.templates["order_completed"] = &ports.EmailTemplate{
		Name:    "order_completed",
		Subject: "Pedido Entregado - {{.TrackingNumber}}",
		HTMLBody: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Pedido Entregado</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #27ae60; border-bottom: 2px solid #27ae60; padding-bottom: 10px;">
¡Pedido Entregado Exitosamente!
        </h2>
        
        <p>Hola <strong>{{.CustomerName}}</strong>,</p>
        
        <p>¡Excelente! Tu pedido ha sido entregado exitosamente.</p>
        
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3>Resumen del Pedido:</h3>
            <ul style="list-style: none; padding: 0;">
                <li><strong>Número de Seguimiento:</strong> {{.TrackingNumber}}</li>
                <li><strong>Estado:</strong> Entregado</li>
                <li><strong>Entregado en:</strong> {{.DeliveryAddress}}</li>
            </ul>
        </div>
        
        <div style="background-color: #e8f5e8; padding: 15px; border-radius: 5px; margin: 20px 0; text-align: center;">
            <h3>¡Gracias por confiar en nosotros!</h3>
            <p>Esperamos que hayas tenido una excelente experiencia. 
               Tu satisfacción es nuestra prioridad.</p>
        </div>
        
        <p style="text-align: center; margin: 30px 0;">
            <a href="mailto:{{.CompanyEmail}}?subject=Feedback - Pedido {{.TrackingNumber}}" 
               style="background-color: #3498db; color: white; padding: 12px 25px; 
                      text-decoration: none; border-radius: 5px; display: inline-block;">
                Enviar Comentarios
            </a>
        </p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        
        <p style="font-size: 12px; color: #666; text-align: center;">
            {{.CompanyName}}<br>
            {{.CompanyPhone}} | {{.CompanyEmail}}<br>
            © {{.CurrentYear}} Todos los derechos reservados
        </p>
    </div>
</body>
</html>`,
		PlainBody: `Hola {{.CustomerName}},

¡Excelente! Tu pedido ha sido entregado exitosamente.

Resumen del Pedido:
- Número de Seguimiento: {{.TrackingNumber}}
- Estado: Entregado
- Entregado en: {{.DeliveryAddress}}

¡Gracias por confiar en nosotros!
Esperamos que hayas tenido una excelente experiencia.

{{.CompanyName}}
{{.CompanyPhone}} | {{.CompanyEmail}}`,
	}
}

func (s *EmailServiceImpl) SendTestEmail(ctx context.Context, to, subject, body, contentType string, cc, bcc []string) (string, error) {
	emailData := &ports.EmailData{
		To:      []string{to},
		CC:      cc,
		BCC:     bcc,
		Subject: subject,
	}

	// Configurar el cuerpo según el tipo de contenido
	if contentType == "html" {
		emailData.Body = body
	} else {
		emailData.PlainText = body
	}

	// Enviar el email
	err := s.SendEmail(ctx, emailData)
	if err != nil {
		return "", err
	}

	// Generar un ID de mensaje simulado
	messageID := fmt.Sprintf("msg_%d", time.Now().Unix())
	return messageID, nil
}

// GetServiceStatus obtiene el estado del servicio de email
func (s *EmailServiceImpl) GetServiceStatus() map[string]interface{} {
	status := map[string]interface{}{
		"smtp_server": s.config.EmailConfig.SMTPHost,
		"smtp_port":   s.config.EmailConfig.SMTPPort,
		"from_email":  s.config.EmailConfig.FromEmail,
		"tls_enabled": s.config.EmailConfig.EnableTLS,
		"last_check":  time.Now().Format(time.RFC3339),
	}

	// Probar conexión SMTP
	if err := s.testSMTPConnection(); err != nil {
		status["connected"] = false
		status["error"] = err.Error()
	} else {
		status["connected"] = true
	}

	return status
}

// testSMTPConnection prueba la conexión SMTP
func (s *EmailServiceImpl) testSMTPConnection() error {
	// Crear un cliente temporal para probar la conexión
	closer, err := s.dialer.Dial()
	if err != nil {
		return err
	}
	defer closer.Close()
	return nil
}

// SendWelcomeTestEmail envía un email de bienvenida de prueba
func (s *EmailServiceImpl) SendWelcomeTestEmail(ctx context.Context, to, name string) (string, error) {
	variables := map[string]interface{}{
		"CustomerName": name,
		"CompanyName":  s.config.EmailConfig.FromName,
		"CompanyEmail": s.config.EmailConfig.FromEmail,
		"CurrentYear":  time.Now().Year(),
		"TestMessage":  "Este es un email de prueba del sistema",
	}

	// Crear un template de bienvenida simple si no existe
	if _, exists := s.templates["welcome_test"]; !exists {
		s.templates["welcome_test"] = &ports.EmailTemplate{
			Name:    "welcome_test",
			Subject: "¡Bienvenido al Sistema de Delivery! - Prueba",
			HTMLBody: `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Bienvenido</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #007bff; color: white; padding: 20px; text-align: center; border-radius: 5px; }
        .content { padding: 20px; background: #f8f9fa; margin: 20px 0; border-radius: 5px; }
        .footer { text-align: center; padding: 10px; font-size: 12px; color: #666; }
        .success { background: #d4edda; border: 1px solid #c3e6cb; color: #155724; padding: 15px; border-radius: 5px; margin: 20px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.CompanyName}}</h1>
            <p>¡Sistema de Delivery en Funcionamiento!</p>
        </div>
        <div class="content">
            <h2>¡Hola {{.CustomerName}}!</h2>
            <p>{{.TestMessage}}</p>
            
            <div class="success">
                <h3>Estado del Sistema:</h3>
                <ul>
                    <li>Configuración de SMTP: OK</li>
                    <li>Envío de emails HTML: OK</li>
                    <li>Templates funcionando: OK</li>
                    <li>Sistema completamente operativo: OK</li>
                </ul>
            </div>
            
            <p><strong>¡Todo está funcionando perfectamente!</strong></p>
            <p>Este email confirma que el servicio de correo electrónico está configurado correctamente y listo para enviar notificaciones de pedidos.</p>
        </div>
        <div class="footer">
            <p>{{.CompanyName}} - {{.CompanyEmail}}</p>
            <p>© {{.CurrentYear}} - Email de prueba automático</p>
        </div>
    </div>
</body>
</html>`,
			PlainBody: `¡Hola {{.CustomerName}}!

{{.TestMessage}}

Estado del Sistema:
Configuración de SMTP: OK
Envío de emails: OK  
Templates funcionando: OK
Sistema completamente operativo: OK

¡Todo está funcionando perfectamente!

Este email confirma que el servicio de correo electrónico está configurado correctamente.

{{.CompanyName}} - {{.CompanyEmail}}
© {{.CurrentYear}} - Email de prueba automático`,
		}
	}

	// Enviar usando el template
	err := s.SendTemplatedEmail(ctx, "welcome_test", []string{to}, variables)
	if err != nil {
		return "", err
	}

	// Generar ID de mensaje
	messageID := fmt.Sprintf("welcome_msg_%d", time.Now().Unix())
	return messageID, nil
}
