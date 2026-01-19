package scanner

import (
	"fmt"
	"sync"
)

// NewScannerFunc is a factory function that creates a scanner instance
type NewScannerFunc func(config Config, timeoutSeconds int) (Scanner, error)

var (
	registry = make(map[string]NewScannerFunc)
	mu       sync.RWMutex
)

// Register registers a scanner backend with the factory
// This should be called from init() functions in backend packages
func Register(name string, fn NewScannerFunc) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = fn
}

// New creates a new scanner instance for the specified backend type
func New(backendType string, config Config, timeoutSeconds int) (Scanner, error) {
	mu.RLock()
	fn, exists := registry[backendType]
	mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unknown scanner backend: %s", backendType)
	}

	return fn(config, timeoutSeconds)
}

// Available returns a list of registered scanner backend names
func Available() []string {
	mu.RLock()
	defer mu.RUnlock()

	backends := make([]string, 0, len(registry))
	for name := range registry {
		backends = append(backends, name)
	}
	return backends
}
