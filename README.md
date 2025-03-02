# Delivery System Backend

Un sistema de gestión de entregas construido con arquitectura hexagonal y Domain-Driven Design en Go.

<img src="assets/logo.png" alt="Delivery System Logo" width="400">

## 📋 Descripción

Este sistema de delivery representa una solución completa diseñada para gestionar el ciclo completo de entregas, desde la recepción de pedidos hasta su entrega final. Su arquitectura está enfocada en la escalabilidad, mantenibilidad y claridad del código.

### Características Principales

- Procesamiento concurrente de grandes volúmenes de pedidos
- Gestión de asignaciones de repartidores en tiempo real
- Sistema de tracking en tiempo real
- Gestión de almacenes y zonas
- Procesamiento de pagos
- Sistema de notificaciones
- Análisis y reportes

## 🏗️ Arquitectura

El sistema está construido siguiendo tres pilares fundamentales:

1. **Domain-Driven Design (DDD)**: El dominio del negocio es el corazón del sistema. Cada concepto de negocio se refleja claramente en el código.

2. **Arquitectura Hexagonal (Ports & Adapters)**: Mantiene el dominio de negocio aislado de las preocupaciones técnicas, facilitando:
    - Independencia de frameworks
    - Testabilidad
    - Mantenibilidad
    - Flexibilidad para cambios

3. **Principios SOLID y Clean Code**: Base fundamental del desarrollo, actuando como guías para asegurar la calidad y mantenibilidad del código.

## 🔧 Tecnologías Utilizadas

- **Go** (v1.23.2+)
- **MySQL/MariaDB** - Para persistencia principal (Se soportan otros motores)
- **Redis** - Para caché y sesiones
- **JWT** - Para autenticación
- **Docker** - Para contenerización
- **Go Modules** y **Go Workspaces** - Para gestión de dependencias

## 📂 Estructura del Proyecto

El proyecto sigue la estructura estándar de Go con adaptaciones para arquitectura hexagonal:

```
.
├── api/                  # Definiciones de API (Swagger)
├── cmd/                  # Punto de entrada de la aplicación
├── config/               # Configuraciones y variables de entorno
├── docs/                 # Documentación
├── githooks/             # Hooks de Git
├── internal/             # Código privado de la aplicación
│   ├── application/      # Casos de uso
│   ├── bootstrap/        # Inicialización y contenedores DI
│   ├── domain/           # Modelos y reglas de negocio
│   └── infrastructure/   # Implementaciones técnicas
├── pkg/                  # Código compartido
│   └── shared/           # Utilidades compartidas (logs, errores)
├── scripts/              # Scripts de utilidad
├── test/                 # Pruebas
│   ├── integration/      # Pruebas de integración
│   ├── mocks/            # Mocks para pruebas
│   ├── performance/      # Pruebas de rendimiento
│   └── unit/             # Pruebas unitarias
├── .env                  # Variables de entorno (local)
├── .env.example          # Plantilla de variables de entorno
├── .gitattributes        # Configuración de Git
├── .gitignore            # Archivos ignorados por Git
├── go.work               # Configuración de workspace Go
├── Makefile              # Comandos de automatización
└── README.md             # Este archivo
```

## ⚙️ Prerrequisitos

- Go 1.23.2 o superior
- MySQL/MariaDB (u otro motor de base de datos)
- Redis
- Docker y Docker Compose (opcional, para desarrollo)
- Make (opcional, para usar comandos del Makefile)

## 🚀 Instalación y Configuración

### Clonación del Repositorio

```bash
git clone https://github.com/MarlonG1/delivery-backend
cd delivery-system
```

### Configuración de Variables de Entorno

```bash
cp .env.example .env
# Edita el archivo .env con tus configuraciones
```

### Configuración de la Base de Datos

El sistema utiliza MySQL como base de datos principal:

```bash
# Crear la base de datos (desde MySQL CLI)
mysql -u root -p
> CREATE DATABASE delivery_db;
> exit

# Alternativa: usando Docker Compose
docker-compose up -d mysql redis
```

### Instalación de Dependencias

```bash
# Inicializa el workspace de Go
go work init

# Añade los módulos al workspace
go work use ./cmd
go work use ./internal
go work use ./pkg
# ... y otros módulos según sea necesario

# Verifica la configuración
go work sync
```

## ▶️ Ejecución

### Desarrollo Local

```bash
# Usando Go directamente
go run cmd/main.go

# Usando Make
make run
```

### Usando Docker

```bash
# Construir imagen
docker build -t delivery-system .

# Ejecutar contenedor
docker run -p 7319:7319 --env-file .env delivery-system
```

## 🧪 Pruebas

```bash
# Ejecutar todas las pruebas
make test

# Ejecutar pruebas unitarias
make test-unit

# Ejecutar pruebas de integración
make test-integration

# Ejecutar pruebas con cobertura
make test-coverage
```

## 📚 Documentación API

La documentación de la API está disponible en formato Swagger:

```bash
# Iniciar servidor de documentación
make swagger-ui

# La documentación estará disponible en:
# http://localhost:8080/swagger/index.html
```

## 📝 Principales Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | /api/v1/auth/login | Iniciar sesión |
| GET | /api/v1/auth/logout | Cerrar sesión |
| GET | /api/v1/users/profile | Obtener perfil de usuario |
| POST | /api/v1/orders | Crear un nuevo pedido |

## 🤝 Contribución

1. Haz un fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/amazing-feature`)
3. Realiza tus cambios y haz commit (`git commit -m 'feat: add amazing feature'`)
4. Haz push a la rama (`git push origin feature/amazing-feature`)
5. Abre un Pull Request

### Convenciones de Commits

Utilizamos una variante de Conventional Commits:

```
<tipo>(<alcance>): <descripción>
```

Tipos: `feat`, `fix`, `refactor`, `docs`, `chore`, etc.

---

Desarrollado con ❤️ por Marlon Isaac Hernández García
```