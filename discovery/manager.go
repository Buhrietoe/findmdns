package discovery

import (
	"fmt"
	"sync"
	"time"

	"github.com/Buhrietoe/findmdns/models"
	"github.com/hashicorp/mdns"
)

// Manager handles mDNS discovery with caching
type Manager struct {
	mu        sync.RWMutex
	devices   map[string]*models.Device
	lastScan  time.Time
}

// NewManager creates a new discovery manager
func NewManager() *Manager {
	return &Manager{
		devices: make(map[string]*models.Device),
	}
}

// Scan performs a new discovery scan
func (m *Manager) Scan() error {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	entries := []*mdns.ServiceEntry{}

	go func() {
		for entry := range entriesCh {
			entries = append(entries, entry)
		}
	}()

	params := mdns.QueryParam{
		Service: "_services._dns-sd._udp",
		Timeout: time.Second * 5,
		Entries: entriesCh,
	}

	mdns.Query(&params)
	close(entriesCh)

	m.updateDevices(entries)
	m.lastScan = time.Now()

	return nil
}

// GetDevices returns all discovered devices
func (m *Manager) GetDevices() []*models.Device {
	m.mu.RLock()
	defer m.mu.RUnlock()

	devices := make([]*models.Device, 0, len(m.devices))
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices
}

// GetDevice returns a specific device by ID
func (m *Manager) GetDevice(id string) (*models.Device, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.devices[id]
	return device, exists
}

// GetLastScan returns the time of the last scan
func (m *Manager) GetLastScan() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastScan
}

// updateDevices updates the internal device cache
func (m *Manager) updateDevices(entries []*mdns.ServiceEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear existing devices
	m.devices = make(map[string]*models.Device)

	// Process new entries
	for _, entry := range entries {
		device := m.entryToDevice(entry)
		m.devices[device.ID] = device
	}
}

// entryToDevice converts an mdns.ServiceEntry to a Device model
func (m *Manager) entryToDevice(entry *mdns.ServiceEntry) *models.Device {
	// Generate a unique ID for the device
	id := fmt.Sprintf("%s-%s-%d", entry.Name, entry.Host, entry.Port)

	// Convert IP addresses to strings
	addresses := []string{}
	
	// Add IPv4 addresses
	if entry.AddrV4 != nil {
		addresses = append(addresses, entry.AddrV4.String())
	}
	
	// Add IPv6 addresses if available  
	if entry.AddrV6 != nil {
		addresses = append(addresses, entry.AddrV6.String())
	}

	device := &models.Device{
		ID:        id,
		Name:      entry.Name,
		Host:      entry.Host,
		Port:      entry.Port,
		Service:   entry.Name, // Using name as service since we're querying _services._dns-sd._udp
		TXT:       make(map[string]string), // TXT records would be in InfoFields if available
		Addresses: addresses,
		LastSeen:  time.Now(),
	}

	// Process TXT records if they exist
	if entry.InfoFields != nil {
		device.TXT = make(map[string]string)
		for _, field := range entry.InfoFields {
			// Parse key=value format from TXT records
			if equalsIndex := indexOf(field, "="); equalsIndex != -1 {
				key := field[:equalsIndex]
				value := field[equalsIndex+1:]
				device.TXT[key] = value
			}
		}
	}

	return device
}

// Helper function to find index of a character in a string
func indexOf(s, substr string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == substr[0] {
			if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
				return i
			}
		}
	}
	return -1
}