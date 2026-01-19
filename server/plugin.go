package main

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/mattermost/mattermost-plugin-antivirus/server/scanner"
	// Import scanner backends to trigger registration
	_ "github.com/mattermost/mattermost-plugin-antivirus/server/scanner/clamav"
)

type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex
	configuration     *configuration

	scannerLock sync.RWMutex
	scanner     scanner.Scanner
}

// OnActivate is called when the plugin is activated
func (p *Plugin) OnActivate() error {
	if err := p.initializeScanner(); err != nil {
		return err
	}

	// Perform health check
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s := p.getScanner()
	if s != nil {
		if err := s.HealthCheck(ctx); err != nil {
			p.API.LogWarn("Scanner health check failed on activation", "error", err.Error())
			// Don't fail activation, but log warning
		} else {
			p.API.LogInfo("Scanner health check passed", "backend", s.Name())
		}
	}

	return nil
}

// OnDeactivate is called when the plugin is deactivated
func (p *Plugin) OnDeactivate() error {
	p.scannerLock.Lock()
	defer p.scannerLock.Unlock()

	if p.scanner != nil {
		return p.scanner.Close()
	}
	return nil
}

// initializeScanner creates and initializes the scanner based on configuration
func (p *Plugin) initializeScanner() error {
	config := p.getConfiguration()

	// Get scanner-specific config
	scannerConfig, err := config.GetScannerConfig()
	if err != nil {
		return err
	}

	// Create scanner
	newScanner, err := scanner.New(
		config.BackendType,
		scannerConfig,
		config.ScanTimeoutSeconds,
	)
	if err != nil {
		return err
	}

	// Replace old scanner
	p.scannerLock.Lock()
	oldScanner := p.scanner
	p.scanner = newScanner
	p.scannerLock.Unlock()

	// Clean up old scanner
	if oldScanner != nil {
		go func() {
			if err := oldScanner.Close(); err != nil {
				p.API.LogError("Error closing old scanner", "error", err.Error())
			}
		}()
	}

	p.API.LogInfo("Initialized scanner", "backend", config.BackendType)
	return nil
}

// getScanner returns the current scanner instance (thread-safe)
func (p *Plugin) getScanner() scanner.Scanner {
	p.scannerLock.RLock()
	defer p.scannerLock.RUnlock()
	return p.scanner
}

// FileWillBeUploaded is called before a file upload is completed
func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, _ io.Writer) (*model.FileInfo, string) {
	config := p.getConfiguration()
	s := p.getScanner()

	if s == nil {
		p.API.LogError("Scanner not initialized")
		return nil, "Antivirus scanner not configured properly. Contact your administrator."
	}

	// Send toast notification
	if err := p.API.SendToastMessage(info.CreatorId, c.SessionId, "Scanning file...", model.SendToastMessageOptions{
		Position: "bottom-center",
	}); err != nil {
		p.API.LogError("Error sending toast message", "error", err.Error())
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(config.ScanTimeoutSeconds)*time.Second,
	)
	defer cancel()

	// Perform scan
	result, err := s.Scan(ctx, file, info.Name)
	if err != nil {
		p.API.LogError("Scan error",
			"filename", info.Name,
			"user", info.CreatorId,
			"error", err.Error(),
		)
	}

	// Handle result based on status
	if result != nil {
		switch result.Status {
		case scanner.ScanStatusClean:
			return info, ""

		case scanner.ScanStatusInfected:
			p.API.LogWarn("Virus detected",
				"filename", info.Name,
				"user", info.CreatorId,
				"threat", result.Message,
				"raw", result.Raw,
			)
			return nil, "This file contains a virus and cannot be uploaded."

		case scanner.ScanStatusTimeout:
			p.API.LogError("Scan timeout",
				"filename", info.Name,
				"user", info.CreatorId,
			)
			return nil, "File scanning timed out. Please try again or contact your administrator."

		case scanner.ScanStatusError:
			p.API.LogError("Scan error",
				"filename", info.Name,
				"user", info.CreatorId,
				"message", result.Message,
			)
			return nil, "An error occurred while scanning your file. Please contact your administrator."

		default:
			p.API.LogError("Unknown scan status",
				"status", result.Status,
				"filename", info.Name,
			)
			return nil, "An unexpected error occurred. Please contact your administrator."
		}
	}

	// Fallback if result is nil
	p.API.LogError("Scan returned nil result",
		"filename", info.Name,
		"user", info.CreatorId,
	)
	return nil, "An error occurred while scanning your file. Please contact your administrator."
}
