package main

import (
	"testing"

	"github.com/mattermost/mattermost-plugin-antivirus/server/scanner/clamav"
)

func TestMigrateConfiguration(t *testing.T) {
	tests := []struct {
		name           string
		input          *configuration
		expectMigrated bool
		checkFields    func(*testing.T, *configuration)
	}{
		{
			name: "v1.x config without BackendType",
			input: &configuration{
				ConnectionType:     "tcp",
				ClamavHostPort:     "localhost:3310",
				ClamavSocketPath:   "/tmp/clamd.socket",
				ScanTimeoutSeconds: 15,
			},
			expectMigrated: true,
			checkFields: func(t *testing.T, c *configuration) {
				if c.BackendType != "clamav" {
					t.Errorf("Expected BackendType 'clamav', got '%s'", c.BackendType)
				}
				if c.ClamAV == nil {
					t.Fatal("Expected ClamAV config to be set, got nil")
				}
				if c.ClamAV.ConnectionType != "tcp" {
					t.Errorf("Expected ClamAV.ConnectionType 'tcp', got '%s'", c.ClamAV.ConnectionType)
				}
				if c.ClamAV.HostPort != "localhost:3310" {
					t.Errorf("Expected ClamAV.HostPort 'localhost:3310', got '%s'", c.ClamAV.HostPort)
				}
				if c.ScanTimeoutSeconds != 15 {
					t.Errorf("Expected ScanTimeoutSeconds 15, got %d", c.ScanTimeoutSeconds)
				}
			},
		},
		{
			name: "v2.x config with BackendType",
			input: &configuration{
				BackendType:        "clamav",
				ScanTimeoutSeconds: 10,
				ClamAV: &clamav.Config{
					ConnectionType: "tcp",
					HostPort:       "localhost:3310",
				},
			},
			expectMigrated: false,
			checkFields: func(t *testing.T, c *configuration) {
				if c.BackendType != "clamav" {
					t.Errorf("Expected BackendType 'clamav', got '%s'", c.BackendType)
				}
			},
		},
		{
			name: "empty v1.x config gets defaults",
			input: &configuration{
				ScanTimeoutSeconds: 0,
			},
			expectMigrated: true,
			checkFields: func(t *testing.T, c *configuration) {
				if c.BackendType != "clamav" {
					t.Errorf("Expected BackendType 'clamav', got '%s'", c.BackendType)
				}
				if c.ScanTimeoutSeconds != 10 {
					t.Errorf("Expected default ScanTimeoutSeconds 10, got %d", c.ScanTimeoutSeconds)
				}
				if c.ClamAV == nil {
					t.Fatal("Expected ClamAV config to be set, got nil")
				}
				if c.ClamAV.ConnectionType != "tcp" {
					t.Errorf("Expected default ConnectionType 'tcp', got '%s'", c.ClamAV.ConnectionType)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock plugin (we only need the API for logging)
			p := &Plugin{}

			result, migrated := p.migrateConfiguration(tt.input)

			if migrated != tt.expectMigrated {
				t.Errorf("Expected migrated=%v, got %v", tt.expectMigrated, migrated)
			}

			if tt.checkFields != nil {
				tt.checkFields(t, result)
			}
		})
	}
}
