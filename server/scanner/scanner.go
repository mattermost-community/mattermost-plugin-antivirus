package scanner

import (
	"context"
	"io"
)

// ScanStatus represents the result of a virus scan
type ScanStatus string

const (
	// ScanStatusClean indicates the file is clean (no threats detected)
	ScanStatusClean ScanStatus = "clean"
	// ScanStatusInfected indicates a virus or threat was detected
	ScanStatusInfected ScanStatus = "infected"
	// ScanStatusError indicates an error occurred during scanning
	ScanStatusError ScanStatus = "error"
	// ScanStatusTimeout indicates the scan timed out
	ScanStatusTimeout ScanStatus = "timeout"
)

// ScanResult contains the outcome of a scan operation
type ScanResult struct {
	Status  ScanStatus // The scan status
	Message string     // Human-readable message (threat name or error description)
	Raw     string     // Raw response from scanner (for logging/debugging)
}

// Scanner is the interface that all antivirus backends must implement
type Scanner interface {
	// Scan performs a virus scan on the provided reader
	// The context can be used for cancellation and timeout control
	// filename is provided for scanners that may use it in their analysis
	Scan(ctx context.Context, reader io.Reader, filename string) (*ScanResult, error)

	// Name returns the scanner backend name (e.g., "clamav", "icap")
	Name() string

	// HealthCheck verifies the scanner is operational
	// Returns nil if healthy, error otherwise
	HealthCheck(ctx context.Context) error

	// Close releases any resources held by the scanner
	// Should be called when the scanner is no longer needed
	Close() error
}

// Config is the interface for scanner-specific configuration
type Config interface {
	// Validate checks if the configuration is valid
	Validate() error
}
