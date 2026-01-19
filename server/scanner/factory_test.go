package scanner

import (
	"testing"
)

func TestRegister(t *testing.T) {
	// Test that Register adds a backend to the registry
	testFn := func(config Config, timeoutSeconds int) (Scanner, error) {
		return nil, nil
	}

	Register("test-backend", testFn)

	backends := Available()
	found := false
	//nolint:modernize // Simple loop is clearer and more explicit here
	for _, backend := range backends {
		if backend == "test-backend" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected 'test-backend' to be registered")
	}
}

func TestNew_UnknownBackend(t *testing.T) {
	// Test that New returns an error for unknown backend
	_, err := New("completely-unknown-backend-xyz", nil, 10)
	if err == nil {
		t.Error("Expected error for unknown backend, got nil")
	}
}
