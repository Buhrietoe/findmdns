package discovery

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Buhrietoe/findmdns/internal/device"
	"github.com/hashicorp/mdns"
)

type Manager struct {
	mu       sync.RWMutex
	devices  map[string]*device.Device
	lastScan time.Time
	scanMu   sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		devices: make(map[string]*device.Device),
	}
}

func (m *Manager) Scan() error {
	m.scanMu.Lock()
	defer m.scanMu.Unlock()

	entriesCh := make(chan *mdns.ServiceEntry, 4)
	var entries []*mdns.ServiceEntry
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for entry := range entriesCh {
			entries = append(entries, entry)
		}
	}()

	params := mdns.QueryParam{
		Service: "_services._dns-sd._udp",
		Timeout: time.Second * 5,
		Entries: entriesCh,
	}

	err := mdns.Query(&params)
	close(entriesCh)
	wg.Wait()

	m.updateDevices(entries)

	m.mu.Lock()
	m.lastScan = time.Now()
	m.mu.Unlock()

	return err
}

func (m *Manager) GetDevices() []*device.Device {
	m.mu.RLock()
	defer m.mu.RUnlock()

	devices := make([]*device.Device, 0, len(m.devices))
	for _, device := range m.devices {
		devices = append(devices, device)
	}
	return devices
}

func (m *Manager) GetDevice(id string) (*device.Device, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	device, exists := m.devices[id]
	return device, exists
}

func (m *Manager) GetLastScan() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastScan
}

func (m *Manager) updateDevices(entries []*mdns.ServiceEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()

	seen := make(map[string]bool)
	for _, entry := range entries {
		device := m.entryToDevice(entry)
		m.devices[device.ID] = device
		seen[device.ID] = true
	}

	for id, device := range m.devices {
		if !seen[id] {
			device.LastSeen = time.Now()
		}
	}
}

func (m *Manager) entryToDevice(entry *mdns.ServiceEntry) *device.Device {
	idSource := fmt.Sprintf("%s-%s-%d", entry.Name, entry.Host, entry.Port)
	hash := sha256.Sum256([]byte(idSource))
	id := hex.EncodeToString(hash[:8])

	addresses := []string{}
	if entry.AddrV4 != nil {
		addresses = append(addresses, entry.AddrV4.String())
	}
	if entry.AddrV6 != nil {
		addresses = append(addresses, entry.AddrV6.String())
	}

	device := &device.Device{
		ID:        id,
		Name:      entry.Name,
		Host:      entry.Host,
		Port:      entry.Port,
		TXT:       make(map[string]string),
		Addresses: addresses,
		LastSeen:  time.Now(),
	}

	if entry.InfoFields != nil {
		for _, field := range entry.InfoFields {
			if equalsIndex := strings.Index(field, "="); equalsIndex != -1 {
				key := field[:equalsIndex]
				value := field[equalsIndex+1:]
				device.TXT[key] = value
			}
		}
	}

	return device
}
