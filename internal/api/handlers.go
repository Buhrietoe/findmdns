package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Buhrietoe/findmdns/internal/discovery"
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

type ScanResponse struct {
	Status  string `json:"status"`
	Devices int    `json:"devices"`
	Error   string `json:"error,omitempty"`
}

type LastScanResponse struct {
	LastScan time.Time `json:"last_scan"`
}

func (h *Handler) ScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := h.manager.Scan()
	devices := h.manager.GetDevices()

	resp := ScanResponse{Status: "completed", Devices: len(devices)}
	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetLastScanHandler(w http.ResponseWriter, r *http.Request) {
	lastScan := h.manager.GetLastScan()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LastScanResponse{LastScan: lastScan})
}
