package user

import (
	"context"
	appPorts "github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UsererUseCase struct {
	userService  interfaces.Userer
	rolesService interfaces.Roler
	compService  interfaces.Companyrer
	tokenService appPorts.TokenProvider
}

func NewUserProfileUseCase(userService interfaces.Userer, rolesService interfaces.Roler, compService interfaces.Companyrer, tokenService appPorts.TokenProvider) appPorts.UserUseCase {
	return &UsererUseCase{
		userService:  userService,
		rolesService: rolesService,
		compService:  compService,
		tokenService: tokenService,
	}
}

func (uc *UsererUseCase) GetProfileInfo(ctx context.Context) (*entities.User, error) {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 1. Obtener la información del usuario
	user, err := uc.userService.GetUserInfo(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UsererUseCase) CreateUser(ctx context.Context, user *entities.User) error {
	// 1. Obtener el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Obtener el CompanyID
	companyID, _, err := uc.compService.GetCompanyAndBranchForUser(ctx, claims.UserID)
	if err != nil {
		return err
	}
	user.CompanyID = companyID

	// 3. Obtener el ID de los roles y asignarlos al usuario
	var roles []entities.Role
	for _, reqRole := range user.Roles {
		dbRol, err := uc.rolesService.GetRoleByIDOrName(ctx, strings.ToUpper(reqRole.Role.Name))
		if err != nil {
			return err
		}

		roles = append(roles, *dbRol)
		user.Roles = nil
	}

	// 4. Crear el usuario
	err = uc.userService.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	// 5. Asignar roles al usuario
	for _, role := range roles {
		err = uc.userService.AssignRoleToUser(ctx, user.ID, role.ID, claims.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *UsererUseCase) UpdateUser(ctx context.Context, userID string, user *entities.User) error {
	// 1. Obtener el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Obtener el ID de los roles y asignarlos al usuario
	var roles []entities.Role
	for _, reqRole := range user.Roles {
		dbRol, err := uc.rolesService.GetRoleByIDOrName(ctx, strings.ToUpper(reqRole.Role.Name))
		if err != nil {
			return err
		}

		roles = append(roles, *dbRol)
		user.Roles = nil
	}

	// 3. Actualizar el usuario
	err := uc.userService.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	// 4. Asignar roles al usuario
	err = uc.userService.UpdateRolesToUser(ctx, userID, claims.UserID, roles)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) DeleteUser(ctx context.Context, userID string) error {
	// 1. Eliminar el usuario
	err := uc.userService.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Obtener la información del usuario
	user, err := uc.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UsererUseCase) RecoverUser(ctx context.Context, id string) error {
	// 1. Recuperar el usuario
	err := uc.userService.RecoverUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) ActivateOrDeactivateUser(ctx context.Context, userID string, active bool) error {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Activar o desactivar el usuario
	err := uc.userService.ActivateOrDeactivateUser(ctx, userID, claims.UserID, active)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) AssignRoleToUser(ctx context.Context, userID, param string) error {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Verificar si el rol existe
	exist, err := uc.rolesService.IsRoleExist(ctx, param)
	if err != nil {
		return err
	}
	if !exist {
		return err
	}

	// 3. Obtener el rol
	role, err := uc.rolesService.GetRoleByIDOrName(ctx, strings.ToUpper(param))
	if err != nil {
		return err
	}

	// 4. Asignar rol al usuario
	err = uc.userService.AssignRoleToUser(ctx, userID, role.ID, claims.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) GetRoles(ctx context.Context) ([]entities.Role, error) {
	// 1. Obtener los roles
	roles, err := uc.rolesService.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (uc *UsererUseCase) GetAllUsers(ctx context.Context, request *http.Request) ([]entities.User, *entities.UserQueryParams, int64, error) {
	// 1. Parsear los parámetros de consulta
	params := uc.parseOrderQueryParams(request)

	// 2. Obtener el ID de los claims del contexto
	claims := request.Context().Value("claims").(*auth.AuthClaims)

	// 3. Obtener el CompanyID
	companyID, _, err := uc.compService.GetCompanyAndBranchForUser(request.Context(), claims.UserID)
	if err != nil {
		return nil, nil, 0, err
	}

	// 4. Obtener los usuarios
	users, total, err := uc.userService.GetAllUsers(ctx, companyID, params)
	if err != nil {
		return nil, nil, 0, err
	}

	return users, params, total, nil
}

func (uc *UsererUseCase) UnassignRole(ctx context.Context, userID, param string) error {
	// 1. Verificar si el rol existe
	role, err := uc.rolesService.GetRoleByIDOrName(ctx, strings.ToUpper(param))
	if err != nil {
		return err
	}

	// 2. Desasignar rol al usuario
	err = uc.userService.UnassignRole(ctx, userID, role.ID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) CleanAllSessions(ctx context.Context, userID string) error {
	// 1. Obtener el usuario por ID
	user, err := uc.userService.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Eliminar de la cache las sesiones del usuario
	if user.Sessions != nil {
		for _, session := range user.Sessions {
			// Me quede arreglando este PANIC
			err = uc.tokenService.RevokeToken(session.Token)
			if err != nil {
				return err
			}
		}
	}

	// 3. Limpiar todas las sesiones del usuario de la DB
	err = uc.userService.CleanAllSessions(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error) {
	// 1. Obtener los roles del usuario
	roles, err := uc.userService.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// parseOrderQueryParams extrae los parámetros de consulta de la request
func (uc *UsererUseCase) parseOrderQueryParams(r *http.Request) *entities.UserQueryParams {
	params := &entities.UserQueryParams{}

	// Filtros
	params.Phone = r.URL.Query().Get("phone")
	params.Email = r.URL.Query().Get("email")
	params.Name = r.URL.Query().Get("name")

	// Opción para incluir pedidos eliminados
	includeDeletedStr := r.URL.Query().Get("include_deleted")
	params.IncludeDeleted = includeDeletedStr == "true" || includeDeletedStr == "1"

	statusStr := r.URL.Query().Get("status")
	params.Status = statusStr == "true" || statusStr == "1"

	// Fechas
	if creationDateStr := r.URL.Query().Get("creation_date"); creationDateStr != "" {
		creationDate, err := time.Parse(time.RFC3339, creationDateStr)
		if err == nil {
			params.CreationDate = &creationDate
		}
	}

	// Paginación
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		} else {
			params.Page = 1 // Default
		}
	} else {
		params.Page = 1 // Default
	}

	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			params.PageSize = pageSize
		} else {
			params.PageSize = 10 // Default
		}
	} else {
		params.PageSize = 10 // Default
	}

	// Ordenamiento
	params.SortBy = r.URL.Query().Get("sort_by")
	params.SortDirection = r.URL.Query().Get("sort_direction")
	if params.SortDirection != "asc" && params.SortDirection != "desc" {
		params.SortDirection = "desc" // Default
	}

	return params
}
