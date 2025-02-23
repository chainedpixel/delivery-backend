package mocks

//go:generate mockgen -source=../../internal/domain/delivery/ports/role_repository_port.go -destination=./role_repository_mock.go -package=mocks
//go:generate mockgen -source=../../internal/domain/delivery/ports/user_repository_port.go -destination=./user_repository_mock.go -package=mocks
//go:generate mockgen -source=../../internal/domain/delivery/ports/user_service_port.go -destination=./user_service_mock.go -package=mocks
//go:generate mockgen -source=../../internal/application/ports/auth_port.go -destination=./auth_port_mock.go -package=mocks
//go:generate mockgen -source=../../internal/application/ports/redis_cache_port.go -destination=./cache_port_mock.go -package=mocks
