package middleware

import (
	"bufio"
	"fmt"
	"infrastructure/api/responser"
	"net"
	"net/http"
	"runtime/debug"
	"shared/logs"
)

type ErrorMiddleware struct {
	responseWriter *responser.ResponseWriter
}

// NewErrorMiddleware crea una nueva instancia de ErrorMiddleware
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{
		responseWriter: responser.NewResponseWriter(),
	}
}

// Handler es un middleware que captura los errores de pánico y los errores de cliente.
func (m *ErrorMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := string(debug.Stack())
				logs.Error("Recovered from panic", map[string]interface{}{
					"error":      err,
					"stackTrace": stackTrace,
					"path":       r.URL.Path,
					"method":     r.Method,
				})

				m.responseWriter.Error(
					w,
					http.StatusInternalServerError,
					"An unexpected error occurred",
					[]string{fmt.Sprintf("%v", err)},
				)
			}
		}()

		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)

		if sw.status >= 400 && !sw.written {
			var msg string
			switch sw.status {
			case http.StatusNotFound:
				msg = "Resource not found"
			case http.StatusMethodNotAllowed:
				msg = "Method not allowed"
			case http.StatusBadRequest:
				msg = "Bad request"
			default:

			}

			logs.Warn("Client error response", map[string]interface{}{
				"status":  sw.status,
				"message": msg,
				"path":    r.URL.Path,
				"method":  r.Method,
			})
			m.responseWriter.Error(w, sw.status, msg, nil)
		} else if sw.status >= 500 {
			logs.Error("Server error response", map[string]interface{}{
				"status": sw.status,
				"path":   r.URL.Path,
				"method": r.Method,
			})
		}
	})
}

// statusWriter es un wrapper para http.ResponseWriter que captura el código de estado
type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

// WriteHeader captura el código de estado
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

// Write captura el código de estado si no se ha escrito antes
func (w *statusWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.written = true
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Hijack implementa el interface http.Hijacker si es necesario, esto sera util cuando implementemos websocket
func (w *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacking not supported")
}
