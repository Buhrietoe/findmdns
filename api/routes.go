package api

import (
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes with proper parameter handling
func SetupRoutes(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	// API endpoints with proper route definitions
	router.HandleFunc("/api/devices", handler.DevicesHandler).Methods("GET")
	router.HandleFunc("/api/device/{id}", handler.DeviceHandler).Methods("GET")
	router.HandleFunc("/api/scan", handler.ScanHandler).Methods("POST")
	router.HandleFunc("/api/last-scan", handler.GetLastScanHandler).Methods("GET")

	return router
}