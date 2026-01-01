package models

import (
	"time"
)

// Device represents an mDNS discovered device
type Device struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Host        string            `json:"host"`
	Port        int               `json:"port"`
	Protocol    string            `json:"protocol"`
	Service     string            `json:"service"`
	TXT         map[string]string `json:"txt,omitempty"`
	Addresses   []string          `json:"addresses"`
	TTL         int               `json:"ttl"`
	LastSeen    time.Time         `json:"last_seen"`
}