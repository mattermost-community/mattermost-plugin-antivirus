package clamav

import (
	"fmt"
)

// Config holds ClamAV-specific configuration
type Config struct {
	ConnectionType string // "tcp" or "unix"
	HostPort       string // For TCP connection (e.g., "localhost:3310")
	SocketPath     string // For Unix socket connection (e.g., "/tmp/clamd.socket")
}

// Validate checks if the ClamAV configuration is valid
func (c *Config) Validate() error {
	if c.ConnectionType != "tcp" && c.ConnectionType != "unix" {
		return fmt.Errorf("ConnectionType must be 'tcp' or 'unix', got '%s'", c.ConnectionType)
	}

	if c.ConnectionType == "tcp" && c.HostPort == "" {
		return fmt.Errorf("HostPort is required for TCP connection")
	}

	if c.ConnectionType == "unix" && c.SocketPath == "" {
		return fmt.Errorf("SocketPath is required for Unix socket connection")
	}

	return nil
}
