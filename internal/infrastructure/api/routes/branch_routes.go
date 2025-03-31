package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterBranchRoutes(router *mux.Router, branchHandler *handlers.BranchHandler) {
	router.HandleFunc("/branches", branchHandler.GetBranches).Methods(http.MethodGet)
	router.HandleFunc("/branches", branchHandler.CreateBranch).Methods(http.MethodPost)
	router.HandleFunc("/branches/{branch_id}", branchHandler.GetBranchByID).Methods(http.MethodGet)
	router.HandleFunc("/branches/{branch_id}", branchHandler.UpdateBranch).Methods(http.MethodPut)
	router.HandleFunc("/branches/reactivate/{branch_id}", branchHandler.ReactivateBranch).Methods(http.MethodGet)
	router.HandleFunc("/branches/deactivate/{branch_id}", branchHandler.DeactivateBranch).Methods(http.MethodGet)

	router.HandleFunc("/branches/zones/{branch_id}", branchHandler.AssignZoneToBranch).Methods(http.MethodPost)
	router.HandleFunc("/branches/available-zones/{branch_id}", branchHandler.GetAvailableZonesForBranch).Methods(http.MethodGet)

	router.HandleFunc("/branches/metrics/{branch_id}", branchHandler.GetBranchMetrics).Methods(http.MethodGet)
}
