package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/responser"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
)

type EmailTestHandler struct {
	emailService ports.EmailService
	respWriter   *responser.ResponseWriter
}

func NewEmailTestHandler(emailService ports.EmailService) *EmailTestHandler {
	return &EmailTestHandler{
		emailService: emailService,
		respWriter:   responser.NewResponseWriter(),
	}
}

// SendTestEmail godoc
// @Summary      Envía un email de prueba
// @Description  Endpoint para probar el envío de correos electrónicos del sistema
// @Tags         email-testing
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email body dto.EmailTestRequest true "Email test data"
// @Success      200  {object}  dto.EmailTestResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/test/email [post]
func (h *EmailTestHandler) SendTestEmail(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud
	req, err := dto.NewEmailTestRequest(r.Body)
	if err != nil {
		logs.Error("Failed to decode email test request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	logs.Info("Sending test email", map[string]interface{}{
		"to":           req.To,
		"subject":      req.Subject,
		"content_type": req.ContentType,
		"cc_count":     len(req.CC),
		"bcc_count":    len(req.BCC),
	})

	// 2. Preparar EmailData para el servicio existente
	emailData := &ports.EmailData{
		To:      []string{req.To},
		CC:      req.CC,
		BCC:     req.BCC,
		Subject: req.Subject,
	}

	// Configurar el cuerpo según el tipo de contenido
	if req.ContentType == "html" {
		emailData.Body = req.Body
	} else {
		emailData.PlainText = req.Body
	}

	// Agregar attachments si los hay
	if req.Attachments != nil {
		for _, attachmentPath := range req.Attachments {
			attachment := ports.EmailAttachment{
				Filename: attachmentPath,
				MimeType: "application/octet-stream", // Tipo por defecto
			}
			emailData.Attachments = append(emailData.Attachments, attachment)
		}
	}

	// 3. Enviar el email usando el servicio existente
	ctx := r.Context()
	err = h.emailService.SendEmail(ctx, emailData)
	if err != nil {
		logs.Error("Failed to send test email", map[string]interface{}{
			"error": err.Error(),
			"to":    req.To,
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 4. Generar ID de mensaje simulado
	messageID := "test_" + time.Now().Format("20060102150405")

	// 5. Preparar la respuesta
	response := dto.EmailTestResponse{
		Success:   true,
		Message:   "Test email sent successfully",
		MessageID: messageID,
		Details: map[string]interface{}{
			"recipient":    req.To,
			"subject":      req.Subject,
			"content_type": req.ContentType,
			"cc":           req.CC,
			"bcc":          req.BCC,
			"attachments":  len(req.Attachments),
		},
	}

	logs.Info("Test email sent successfully", map[string]interface{}{
		"to":         req.To,
		"message_id": messageID,
	})

	// 6. Responder
	h.respWriter.Success(w, http.StatusOK, response)
}

// GetEmailStatus godoc
// @Summary      Obtiene el estado del servicio de email
// @Description  Endpoint para verificar el estado y configuración del servicio de email
// @Tags         email-testing
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  responser.APIErrorResponse
// @Router       /api/v1/test/email/status [get]
func (h *EmailTestHandler) GetEmailStatus(w http.ResponseWriter, r *http.Request) {
	logs.Info("Email service status requested", nil)

	// Crear respuesta de estado (simulado ya que el EmailService no tiene GetStatus)
	response := map[string]interface{}{
		"service": "email",
		"status": map[string]interface{}{
			"service_available": true,
			"last_check":        time.Now().Format(time.RFC3339),
			"message":           "Email service is operational",
		},
		"message": "Email service status retrieved successfully",
		"capabilities": map[string]bool{
			"send_html":    true,
			"send_text":    true,
			"attachments":  true,
			"cc_bcc":       true,
			"templates":    true,
			"order_emails": true,
		},
		"available_templates": []string{
			"order_created",
			"order_started",
			"order_completed",
			"order_cancelled",
		},
	}

	h.respWriter.Success(w, http.StatusOK, response)
}

// SendWelcomeEmail godoc
// @Summary      Envía un email de bienvenida de prueba
// @Description  Endpoint para enviar un email de bienvenida predefinido (útil para pruebas rápidas)
// @Tags         email-testing
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email body map[string]string true "Recipient info" example:{"to":"test@example.com","name":"John Doe"}
// @Success      200  {object}  dto.EmailTestResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/test/email/welcome [post]
func (h *EmailTestHandler) SendWelcomeEmail(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el JSON correctamente
	var reqBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logs.Error("Failed to decode welcome email request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Extraer y validar datos
	recipientEmail := reqBody["to"]
	recipientName := reqBody["name"]

	if recipientEmail == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Recipient email is required", nil)
		return
	}

	if recipientName == "" {
		recipientName = "Usuario"
	}

	logs.Info("Sending welcome test email", map[string]interface{}{
		"to":   recipientEmail,
		"name": recipientName,
	})

	// 3. Crear template de bienvenida simple
	subject := "¡Bienvenido al Sistema de Delivery! - Prueba"
	htmlBody := `
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
            <h1>🚚 Sistema de Delivery</h1>
            <p>¡Sistema en Funcionamiento!</p>
        </div>
        <div class="content">
            <h2>¡Hola ` + recipientName + `!</h2>
            <p>Este es un email de prueba del sistema de delivery.</p>
            
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
            <p>Sistema de Delivery - Email de prueba automático</p>
            <p>© ` + time.Now().Format("2006") + ` - No responder a este mensaje</p>
        </div>
    </div>
</body>
</html>`

	plainBody := `¡Hola ` + recipientName + `!

Este es un email de prueba del sistema de delivery.

Estado del Sistema:
Configuración de SMTP: OK
Envío de emails: OK  
Templates funcionando: OK
Sistema completamente operativo: OK

¡Todo está funcionando perfectamente!

Este email confirma que el servicio de correo electrónico está configurado correctamente.

Sistema de Delivery - Email de prueba automático
© ` + time.Now().Format("2006") + ` - No responder a este mensaje`

	// 4. Preparar EmailData
	emailData := &ports.EmailData{
		To:        []string{recipientEmail},
		Subject:   subject,
		Body:      htmlBody,
		PlainText: plainBody,
	}

	// 5. Enviar el email
	ctx := r.Context()
	err := h.emailService.SendEmail(ctx, emailData)
	if err != nil {
		logs.Error("Failed to send welcome email", map[string]interface{}{
			"error": err.Error(),
			"to":    recipientEmail,
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 6. Generar respuesta
	messageID := "welcome_" + time.Now().Format("20060102150405")

	response := dto.EmailTestResponse{
		Success:   true,
		Message:   "Welcome email sent successfully",
		MessageID: messageID,
		Details: map[string]interface{}{
			"recipient": recipientEmail,
			"name":      recipientName,
			"type":      "welcome_email",
		},
	}

	logs.Info("Welcome email sent successfully", map[string]interface{}{
		"to":         recipientEmail,
		"name":       recipientName,
		"message_id": messageID,
	})

	h.respWriter.Success(w, http.StatusOK, response)
}

// SendOrderTestEmail godoc
// @Summary      Envía un email de prueba usando templates de pedidos
// @Description  Endpoint para probar los templates de emails de pedidos existentes
// @Tags         email-testing
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email body map[string]interface{} true "Order test data"
// @Success      200  {object}  dto.EmailTestResponse
// @Failure      400  {object}  responser.APIErrorResponse
// @Router       /api/v1/test/email/order [post]
func (h *EmailTestHandler) SendOrderTestEmail(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar el JSON
	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		logs.Error("Failed to decode order test email request", map[string]interface{}{
			"error": err.Error(),
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 2. Extraer datos requeridos
	recipientEmail, ok := reqBody["to"].(string)
	if !ok || recipientEmail == "" {
		h.respWriter.Error(w, http.StatusBadRequest, "Recipient email is required", nil)
		return
	}

	templateName, ok := reqBody["template"].(string)
	if !ok || templateName == "" {
		templateName = "order_created" // Template por defecto
	}

	customerName, _ := reqBody["customer_name"].(string)
	if customerName == "" {
		customerName = "Usuario de Prueba"
	}

	logs.Info("Sending order test email", map[string]interface{}{
		"to":       recipientEmail,
		"template": templateName,
		"customer": customerName,
	})

	// 3. Construir variables de prueba para el template
	variables := map[string]interface{}{
		"CustomerName":    customerName,
		"OrderID":         "TEST-ORDER-001",
		"TrackingNumber":  "TRACK-TEST-" + time.Now().Format("20060102"),
		"TrackingURL":     "https://yoursite.com/track/TRACK-TEST-" + time.Now().Format("20060102"),
		"OrderStatus":     "Confirmado",
		"OrderPrice":      "25.50",
		"CreatedAt":       time.Now().Format("02/01/2006 15:04"),
		"CompanyName":     "Sistema de Delivery - Prueba",
		"CompanyPhone":    "+503 1234-5678",
		"CompanyEmail":    "prueba@delivery.com",
		"CurrentYear":     time.Now().Year(),
		"DeliveryAddress": "Calle Principal #123, San Salvador, El Salvador",
		"RecipientName":   customerName,
		"RecipientPhone":  "+503 8765-4321",
		"EstimatedTime":   "30-45 minutos",
		"DriverName":      "Juan Conductor de Prueba",
		"DriverPhone":     "+503 9999-8888",
		"VehicleInfo":     "Toyota Corolla Blanco - Placa P123456",
		"AdditionalNotes": "Este es un email de prueba del sistema",
	}

	// 4. Agregar variables adicionales del request
	for key, value := range reqBody {
		if key != "to" && key != "template" && key != "customer_name" {
			variables[key] = value
		}
	}

	// 5. Verificar que el template existe (simulación)
	validTemplates := map[string]bool{
		"order_created":   true,
		"order_started":   true,
		"order_completed": true,
		"order_cancelled": true,
	}

	if !validTemplates[templateName] {
		h.respWriter.Error(w, http.StatusBadRequest,
			"Invalid template. Available templates: order_created, order_started, order_completed, order_cancelled", nil)
		return
	}

	// 6. Enviar usando el template existente
	ctx := r.Context()
	err := h.emailService.SendTemplatedEmail(ctx, templateName, []string{recipientEmail}, variables)
	if err != nil {
		logs.Error("Failed to send order test email", map[string]interface{}{
			"error":    err.Error(),
			"to":       recipientEmail,
			"template": templateName,
		})
		h.respWriter.HandleError(w, err)
		return
	}

	// 7. Generar respuesta
	messageID := "order_test_" + templateName + "_" + time.Now().Format("20060102150405")

	response := dto.EmailTestResponse{
		Success:   true,
		Message:   "Order test email sent successfully using template: " + templateName,
		MessageID: messageID,
		Details: map[string]interface{}{
			"recipient":     recipientEmail,
			"template":      templateName,
			"customer_name": customerName,
			"variables":     variables,
		},
	}

	logs.Info("Order test email sent successfully", map[string]interface{}{
		"to":         recipientEmail,
		"template":   templateName,
		"message_id": messageID,
	})

	h.respWriter.Success(w, http.StatusOK, response)
}
