FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar herramientas necesarias
RUN apk add --no-cache git

# Copiar proyecto
COPY . .

WORKDIR /app
RUN go build -o delivery-app cmd/main.go

# Verificar dónde quedó el binario
RUN echo "=== Buscando binario compilado ===" && \
    find / -name "delivery-app" -type f 2>/dev/null

# Imagen final
FROM alpine:latest

WORKDIR /app

# Copiar binario compilado
COPY --from=builder /app/delivery-app* /app/ || true
COPY --from=builder /app/cmd/delivery-app* /app/ || true
COPY --from=builder /delivery-app* /app/ || true

RUN ls -la /app/

RUN apk add --no-cache ca-certificates

# Configuración
COPY .env /app/

RUN chmod +x /app/delivery-app || echo "No se encontró el binario para asignar permisos"

EXPOSE 7319

CMD ["sh", "-c", "ls -la /app && exec /app/delivery-app"]