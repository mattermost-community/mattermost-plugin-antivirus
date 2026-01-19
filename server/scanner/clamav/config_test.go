package clamav

import (
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid TCP config",
			config: &Config{
				ConnectionType: "tcp",
				HostPort:       "localhost:3310",
			},
			wantErr: false,
		},
		{
			name: "valid Unix config",
			config: &Config{
				ConnectionType: "unix",
				SocketPath:     "/tmp/clamd.socket",
			},
			wantErr: false,
		},
		{
			name: "invalid connection type",
			config: &Config{
				ConnectionType: "invalid",
				HostPort:       "localhost:3310",
			},
			wantErr: true,
		},
		{
			name: "TCP without HostPort",
			config: &Config{
				ConnectionType: "tcp",
			},
			wantErr: true,
		},
		{
			name: "Unix without SocketPath",
			config: &Config{
				ConnectionType: "unix",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
