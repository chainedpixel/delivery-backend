package dto

// ZoneAssignmentRequest representa la solicitud para asignar una zona a una sucursal
// @Description Solicitud para asignar una zona a una sucursal
type ZoneAssignmentRequest struct {
	// ID de la zona a asignar
	// @required
	ZoneID string `json:"zone_id" example:"f8c3e8d7-b6a5-4d3c-9f1e-0a2b4c6d8e0f" binding:"required"`
}
