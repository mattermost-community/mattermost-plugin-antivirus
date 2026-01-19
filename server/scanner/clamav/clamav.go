package clamav

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/IntelXLabs-LLC/go-clamd"

	"github.com/mattermost/mattermost-plugin-antivirus/server/scanner"
)

func init() {
	// Register ClamAV scanner with the factory
	scanner.Register("clamav", NewScanner)
}

// Scanner implements the scanner.Scanner interface for ClamAV
type Scanner struct {
	config  *Config
	timeout time.Duration
}

// NewScanner creates a new ClamAV scanner instance
func NewScanner(cfg scanner.Config, timeoutSeconds int) (scanner.Scanner, error) {
	config, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("invalid config type for ClamAV scanner")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ClamAV config: %w", err)
	}

	return &Scanner{
		config:  config,
		timeout: time.Duration(timeoutSeconds) * time.Second,
	}, nil
}

// Scan performs a virus scan using ClamAV
func (s *Scanner) Scan(ctx context.Context, reader io.Reader, filename string) (*scanner.ScanResult, error) {
	// Create ClamAV client
	var av *clamd.Clamd
	if s.config.ConnectionType == "tcp" {
		av = clamd.NewClamd("tcp://" + s.config.HostPort)
	} else {
		av = clamd.NewClamd(s.config.SocketPath)
	}

	// Create abort channel for ClamAV
	abortChan := make(chan bool)
	defer close(abortChan)

	// Start scan
	responseChan, err := av.ScanStream(reader, abortChan)
	if err != nil {
		return &scanner.ScanResult{
			Status:  scanner.ScanStatusError,
			Message: "Failed to connect to ClamAV daemon",
			Raw:     err.Error(),
		}, err
	}

	// Wait for result or timeout
	for {
		select {
		case result, ok := <-responseChan:
			if !ok {
				// Channel closed without error, file is clean
				return &scanner.ScanResult{
					Status:  scanner.ScanStatusClean,
					Message: "No threats detected",
					Raw:     "CLEAN",
				}, nil
			}

			// Check scan result
			if result.Status != clamd.RES_OK {
				return &scanner.ScanResult{
					Status:  scanner.ScanStatusInfected,
					Message: result.Description,
					Raw:     result.Raw,
				}, nil
			}
			// Continue reading from channel for additional results
			continue

		case <-ctx.Done():
			// Context canceled or timed out
			return &scanner.ScanResult{
				Status:  scanner.ScanStatusTimeout,
				Message: "Scan timed out or was canceled",
				Raw:     "TIMEOUT",
			}, ctx.Err()

		case <-time.After(s.timeout):
			// Explicit timeout
			return &scanner.ScanResult{
				Status:  scanner.ScanStatusTimeout,
				Message: "Scan timed out",
				Raw:     "TIMEOUT",
			}, fmt.Errorf("scan timeout after %v", s.timeout)
		}
	}
}

// Name returns the scanner backend name
func (s *Scanner) Name() string {
	return "clamav"
}

// HealthCheck verifies that ClamAV is reachable and operational
func (s *Scanner) HealthCheck(ctx context.Context) error {
	var av *clamd.Clamd
	if s.config.ConnectionType == "tcp" {
		av = clamd.NewClamd("tcp://" + s.config.HostPort)
	} else {
		av = clamd.NewClamd(s.config.SocketPath)
	}

	err := av.Ping()
	if err != nil {
		return fmt.Errorf("ClamAV health check failed: %w", err)
	}
	return nil
}

// Close releases any resources held by the scanner
func (s *Scanner) Close() error {
	// ClamAV client doesn't maintain persistent connections, no cleanup needed
	return nil
}
