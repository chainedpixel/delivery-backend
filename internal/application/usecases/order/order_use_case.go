package order

import (
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/entities"
	"infrastructure/api/dto"
	"net/http"
	"shared/mappers/request_mapper"
	"strconv"
	"time"
)

type OrderUseCase struct {
	orderService   interfaces.Orderer
	companyService interfaces.Companyrer
}

func NewOrderUseCase(orderService interfaces.Orderer, companyService interfaces.Companyrer) *OrderUseCase {
	return &OrderUseCase{
		orderService:   orderService,
		companyService: companyService,
	}
}

// CreateOrder crea un nuevo pedido
func (uc *OrderUseCase) CreateOrder(ctx context.Context, authUserID string, reqOrder *dto.OrderCreateRequest) error {
	//1. Obtener la dirección de la empresa según el ID
	companyAddress, err := uc.companyService.GetAddressByID(ctx, reqOrder.CompanyPickUpID)
	if err != nil {
		return err
	}

	// 2. Usar el mapper para convertir el dto a entidad
	order, err := request_mapper.OrderRequestToOrder(reqOrder, companyAddress)
	if err != nil {
		return err
	}

	// 3. Obtener el branch y company ID del usuario
	order.CompanyID, order.BranchID, err = uc.companyService.GetCompanyAndBranchForUser(ctx, authUserID)
	if err != nil {
		return err
	}

	//3. Create order
	err = uc.orderService.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByID obtiene un pedido por su ID
func (uc *OrderUseCase) GetOrderByID(ctx context.Context, orderID string) (*entities.Order, error) {
	order, err := uc.orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// ChangeStatus cambia el estado de un pedido
func (uc *OrderUseCase) ChangeStatus(ctx context.Context, id, status string) error {
	err := uc.orderService.ChangeStatus(ctx, id, status)
	if err != nil {
		return err
	}

	return nil
}

// GetOrdersByCompany obtiene los pedidos de una empresa
func (uc *OrderUseCase) GetOrdersByCompany(ctx context.Context, userID string, request *http.Request) ([]entities.Order, *entities.OrderQueryParams, int64, error) {
	// 1. Parsear los parámetros de consulta
	params := uc.parseOrderQueryParams(request)

	// 2. Obtener el ID de la empresa por el ID del usuario
	companyID, _, err := uc.companyService.GetCompanyAndBranchForUser(ctx, userID)

	// 2. Obtener los pedidos
	orders, total, err := uc.orderService.GetOrdersByCompany(ctx, companyID, params)
	if err != nil {
		return nil, nil, 0, err
	}

	return orders, params, total, nil
}

// parseOrderQueryParams extrae los parámetros de consulta de la request
func (uc *OrderUseCase) parseOrderQueryParams(r *http.Request) *entities.OrderQueryParams {
	params := &entities.OrderQueryParams{}

	// Filtros
	params.Status = r.URL.Query().Get("status")
	params.Location = r.URL.Query().Get("location")
	params.TrackingNumber = r.URL.Query().Get("tracking_number")

	// Fechas
	if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err == nil {
			params.StartDate = &startDate
		}
	}

	if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err == nil {
			params.EndDate = &endDate
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
