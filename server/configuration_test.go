package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfiguration_FromMap(t *testing.T) {
	t.Run("populates all fields from map", func(t *testing.T) {
		m := map[string]interface{}{
			"ClamavHostPort":       "localhost:3310",
			"ScanTimeoutSeconds":   float64(30),
			"ConnectionType":       "tcp",
			"ClamavSocketPath":     "/tmp/clamd.socket",
			"ToastMessageScanning": "Scanning...",
			"ToastMessageSuccess":  "Done!",
		}

		var cfg configuration
		err := cfg.FromMap(m)

		require.NoError(t, err)
		assert.Equal(t, "localhost:3310", cfg.ClamavHostPort)
		assert.Equal(t, 30, cfg.ScanTimeoutSeconds)
		assert.Equal(t, "tcp", cfg.ConnectionType)
		assert.Equal(t, "/tmp/clamd.socket", cfg.ClamavSocketPath)
		assert.Equal(t, "Scanning...", cfg.ToastMessageScanning)
		assert.Equal(t, "Done!", cfg.ToastMessageSuccess)
	})

	t.Run("handles empty map", func(t *testing.T) {
		m := map[string]interface{}{}

		var cfg configuration
		err := cfg.FromMap(m)

		require.NoError(t, err)
		assert.Empty(t, cfg.ClamavHostPort)
		assert.Empty(t, cfg.ToastMessageScanning)
	})

	t.Run("handles nil map", func(t *testing.T) {
		var cfg configuration
		err := cfg.FromMap(nil)

		require.NoError(t, err)
	})
}

func TestConfiguration_ToMap(t *testing.T) {
	t.Run("converts all fields to map", func(t *testing.T) {
		cfg := configuration{
			ClamavHostPort:       "localhost:3310",
			ScanTimeoutSeconds:   30,
			ConnectionType:       "tcp",
			ClamavSocketPath:     "/tmp/clamd.socket",
			ToastMessageScanning: "Scanning...",
			ToastMessageSuccess:  "Done!",
		}

		m, err := cfg.ToMap()

		require.NoError(t, err)
		assert.Equal(t, "localhost:3310", m["clamavhostport"])
		assert.Equal(t, float64(30), m["scantimeoutseconds"])
		assert.Equal(t, "tcp", m["connectiontype"])
		assert.Equal(t, "/tmp/clamd.socket", m["clamavsocketpath"])
		assert.Equal(t, "Scanning...", m["toastmessagescanning"])
		assert.Equal(t, "Done!", m["toastmessagesuccess"])
	})

	t.Run("handles empty configuration", func(t *testing.T) {
		cfg := configuration{}

		m, err := cfg.ToMap()

		require.NoError(t, err)
		assert.NotNil(t, m)
	})
}

func TestConfiguration_Defaults(t *testing.T) {
	t.Run("sets defaults for empty toast messages", func(t *testing.T) {
		cfg := configuration{}

		cfg.Defaults()

		assert.Equal(t, DefaultToastMessageScanning, cfg.ToastMessageScanning)
		assert.Equal(t, DefaultToastMessageSuccess, cfg.ToastMessageSuccess)
	})

	t.Run("trims whitespace from toast messages", func(t *testing.T) {
		cfg := configuration{
			ToastMessageScanning: "  Scanning...  ",
			ToastMessageSuccess:  "  Done!  ",
		}

		cfg.Defaults()

		assert.Equal(t, "Scanning...", cfg.ToastMessageScanning)
		assert.Equal(t, "Done!", cfg.ToastMessageSuccess)
	})

	t.Run("sets defaults for whitespace-only toast messages", func(t *testing.T) {
		cfg := configuration{
			ToastMessageScanning: "   ",
			ToastMessageSuccess:  "\t\n",
		}

		cfg.Defaults()

		assert.Equal(t, DefaultToastMessageScanning, cfg.ToastMessageScanning)
		assert.Equal(t, DefaultToastMessageSuccess, cfg.ToastMessageSuccess)
	})

	t.Run("preserves non-empty toast messages", func(t *testing.T) {
		cfg := configuration{
			ToastMessageScanning: "Custom scanning message",
			ToastMessageSuccess:  "Custom success message",
		}

		cfg.Defaults()

		assert.Equal(t, "Custom scanning message", cfg.ToastMessageScanning)
		assert.Equal(t, "Custom success message", cfg.ToastMessageSuccess)
	})

	t.Run("does not modify other fields", func(t *testing.T) {
		cfg := configuration{
			ClamavHostPort:     "localhost:3310",
			ScanTimeoutSeconds: 30,
			ConnectionType:     "tcp",
			ClamavSocketPath:   "/tmp/clamd.socket",
		}

		cfg.Defaults()

		assert.Equal(t, "localhost:3310", cfg.ClamavHostPort)
		assert.Equal(t, 30, cfg.ScanTimeoutSeconds)
		assert.Equal(t, "tcp", cfg.ConnectionType)
		assert.Equal(t, "/tmp/clamd.socket", cfg.ClamavSocketPath)
	})
}
