package api

import (
	"encoding/json"
	"net/http"

	"github.com/Buhrietoe/findmdns/discovery"
	"github.com/gorilla/mux"
)

// Handler holds the discovery manager
type Handler struct {
	manager *discovery.Manager
}

// NewHandler creates a new API handler
func NewHandler(manager *discovery.Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

// DevicesHandler returns all discovered devices
func (h *Handler) DevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices := h.manager.GetDevices()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// DeviceHandler returns a specific device by ID
func (h *Handler) DeviceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path using gorilla/mux
	vars := mux.Vars(r)
	id := vars["id"]
	
	device, exists := h.manager.GetDevice(id)
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// ScanHandler triggers a new discovery scan
func (h *Handler) ScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	err := h.manager.Scan()
	if err != nil {
		http.Error(w, "Scan failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "scanning started"})
}

// GetLastScanHandler returns the time of the last scan
func (h *Handler) GetLastScanHandler(w http.ResponseWriter, r *http.Request) {
	lastScan := h.manager.GetLastScan()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"last_scan": lastScan,
	})
}