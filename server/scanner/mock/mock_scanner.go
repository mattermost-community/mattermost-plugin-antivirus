package mock

import (
	"context"
	"io"

	"github.com/mattermost/mattermost-plugin-antivirus/server/scanner"
)

// MockScanner is a mock implementation of the scanner.Scanner interface for testing
type MockScanner struct {
	ScanFunc        func(ctx context.Context, reader io.Reader, filename string) (*scanner.ScanResult, error)
	NameFunc        func() string
	HealthCheckFunc func(ctx context.Context) error
	CloseFunc       func() error
}

// Scan calls the mock ScanFunc or returns a clean result by default
func (m *MockScanner) Scan(ctx context.Context, reader io.Reader, filename string) (*scanner.ScanResult, error) {
	if m.ScanFunc != nil {
		return m.ScanFunc(ctx, reader, filename)
	}
	return &scanner.ScanResult{
		Status:  scanner.ScanStatusClean,
		Message: "Mock scanner: file is clean",
		Raw:     "MOCK_CLEAN",
	}, nil
}

// Name calls the mock NameFunc or returns "mock" by default
func (m *MockScanner) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock"
}

// HealthCheck calls the mock HealthCheckFunc or returns nil by default
func (m *MockScanner) HealthCheck(ctx context.Context) error {
	if m.HealthCheckFunc != nil {
		return m.HealthCheckFunc(ctx)
	}
	return nil
}

// Close calls the mock CloseFunc or returns nil by default
func (m *MockScanner) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}
