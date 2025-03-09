package order

import (
	"context"
	"domain/delivery/interfaces"
	"infrastructure/api/dto"
	"shared/mappers/request_mapper"
)

type OrderCreateUseCase struct {
	orderService   interfaces.Orderer
	companyService interfaces.Companyrer
}

func NewOrderUseCase(orderService interfaces.Orderer, companyService interfaces.Companyrer) *OrderCreateUseCase {
	return &OrderCreateUseCase{
		orderService:   orderService,
		companyService: companyService,
	}
}

func (uc *OrderCreateUseCase) CreateOrder(ctx context.Context, reqOrder *dto.OrderCreateRequest) error {
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

	//3. Create order
	err = uc.orderService.CreateOrder(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
