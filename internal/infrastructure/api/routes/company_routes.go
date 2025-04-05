package routes

import (
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterCompanyRoutes(router *mux.Router, companyHandler *handlers.CompanyHandler) {
	router.HandleFunc("/companies/profile", companyHandler.GetCompanyProfile).Methods(http.MethodGet)
	router.HandleFunc("/companies/profile", companyHandler.UpdateCompany).Methods(http.MethodPut)

	router.HandleFunc("/companies/addresses", companyHandler.GetCompanyAddresses).Methods(http.MethodGet)
	router.HandleFunc("/companies/addresses", companyHandler.AddCompanyAddress).Methods(http.MethodPost)
	router.HandleFunc("/companies/addresses/{address_id}", companyHandler.UpdateCompanyAddress).Methods(http.MethodPut)
	router.HandleFunc("/companies/addresses/{address_id}", companyHandler.DeleteCompanyAddress).Methods(http.MethodDelete)

	router.HandleFunc("/companies/deactivate", companyHandler.DeactivateCompany).Methods(http.MethodPost)
	router.HandleFunc("/companies/reactivate", companyHandler.ReactivateCompany).Methods(http.MethodPost)

	router.HandleFunc("/companies/metrics", companyHandler.GetCompanyMetrics).Methods(http.MethodGet)

	router.HandleFunc("/companies", companyHandler.CreateCompany).Methods(http.MethodPost)
	router.HandleFunc("/companies", companyHandler.GetCompanies).Methods(http.MethodGet)
}
