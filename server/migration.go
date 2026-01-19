package main

import (
	"github.com/mattermost/mattermost-plugin-antivirus/server/scanner/clamav"
)

// migrateConfiguration handles backward compatibility for configs from v1.x
// It detects old configs (missing BackendType) and migrates them to v2.x format
func (p *Plugin) migrateConfiguration(config *configuration) (*configuration, bool) {
	migrated := false

	// Check if BackendType is missing (v1.x config)
	if config.BackendType == "" {
		if p.API != nil {
			p.API.LogInfo("Detected v1.x configuration format, migrating to v2.x")
		}

		// Default to ClamAV backend
		config.BackendType = "clamav"
		migrated = true

		// Migrate legacy flat fields to nested ClamAV config
		if config.ClamAV == nil {
			config.ClamAV = &clamav.Config{}
		}

		// Map old fields to new nested structure
		if config.ConnectionType != "" {
			config.ClamAV.ConnectionType = config.ConnectionType
		}
		if config.ClamavHostPort != "" {
			config.ClamAV.HostPort = config.ClamavHostPort
		}
		if config.ClamavSocketPath != "" {
			config.ClamAV.SocketPath = config.ClamavSocketPath
		}

		// Set defaults if not specified
		if config.ClamAV.ConnectionType == "" {
			config.ClamAV.ConnectionType = "tcp"
		}
		if config.ClamAV.HostPort == "" && config.ClamAV.ConnectionType == "tcp" {
			config.ClamAV.HostPort = "localhost:3310"
		}
		if config.ClamAV.SocketPath == "" && config.ClamAV.ConnectionType == "unix" {
			config.ClamAV.SocketPath = "/tmp/clamd.socket"
		}

		// Set default timeout if not specified
		if config.ScanTimeoutSeconds == 0 {
			config.ScanTimeoutSeconds = 10
		}

		if p.API != nil {
			p.API.LogInfo("Successfully migrated configuration to v2.x format",
				"backend_type", config.BackendType,
				"connection_type", config.ClamAV.ConnectionType,
			)
		}
	}

	return config, migrated
}
