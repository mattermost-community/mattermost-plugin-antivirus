package main

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// configuration captures the plugin's external configuration as exposed in the Mattermost server
// configuration, as well as values computed from the configuration. Any public fields will be
// deserialized from the Mattermost server configuration in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// configuration can change at any time, access to the configuration must be synchronized. The
// strategy used in this plugin is to guard a pointer to the configuration, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your configuration struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type configuration struct {
	ClamavHostPort     string `json:"clamavhostport"`
	ScanTimeoutSeconds int    `json:"scantimeoutseconds"`
	ConnectionType     string `json:"connectiontype"`
	ClamavSocketPath   string `json:"clamavsocketpath"`

	// Toast message customization
	ToastMessageScanning string `json:"toastmessagescanning"`
	ToastMessageSuccess  string `json:"toastmessagesuccess"`
}

const (
	DefaultToastMessageScanning = "Scanning file..."
	DefaultToastMessageSuccess  = "File scanned, no threats found"
)

// FromMap populates the configuration from a map[string]interface{}.
func (c *configuration) FromMap(m map[string]any) error {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "failed to marshal plugin configuration")
	}
	if err := json.Unmarshal(jsonBytes, &c); err != nil {
		return errors.Wrap(err, "failed to unmarshal plugin configuration")
	}
	return nil
}

// ToMap converts the configuration to a map[string]interface{}.
func (c *configuration) ToMap() (map[string]any, error) {
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal configuration")
	}
	var m map[string]any
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration")
	}
	return m, nil
}

// Defaults trims string fields and sets default values for empty toast messages.
func (c *configuration) Defaults() {
	c.ToastMessageScanning = strings.TrimSpace(c.ToastMessageScanning)
	c.ToastMessageSuccess = strings.TrimSpace(c.ToastMessageSuccess)

	if c.ToastMessageScanning == "" {
		c.ToastMessageScanning = DefaultToastMessageScanning
	}
	if c.ToastMessageSuccess == "" {
		c.ToastMessageSuccess = DefaultToastMessageSuccess
	}
}

// Clone shallow copies the configuration. Your implementation may require a deep copy if
// your configuration has reference types.
func (c *configuration) Clone() *configuration {
	clone := *c
	return &clone
}

// getConfiguration retrieves the active configuration under lock, making it safe to use
// concurrently. The active configuration may change underneath the client of this method, but
// the struct returned by this API call is considered immutable.
func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

// setConfiguration replaces the active configuration under lock.
//
// Do not call setConfiguration while holding the configurationLock, as sync.Mutex is not
// reentrant. In particular, avoid using the plugin API entirely, as this may in turn trigger a
// hook back into the plugin. If that hook attempts to acquire this lock, a deadlock may occur.
//
// This method panics if setConfiguration is called with the existing configuration. This almost
// certainly means that the configuration was modified without being cloned and may result in
// an unsafe access.
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

// OnConfigurationChange is invoked when configuration changes may have been made.
func (p *Plugin) OnConfigurationChange() error {
	configuration := new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	configuration.Defaults()
	p.setConfiguration(configuration)

	return nil
}

// // ConfigurationWillBeSaved is invoked before saving the configuration to the backing store.
// func (p *Plugin) ConfigurationWillBeSaved(newCfg *model.Config) (*model.Config, error) {
// 	if newCfg == nil || newCfg.PluginSettings.Plugins == nil {
// 		return newCfg, nil
// 	}

// 	pluginConfig, ok := newCfg.PluginSettings.Plugins["antivirus"]
// 	if !ok {
// 		return newCfg, nil
// 	}

// 	// Log the plugin configuration for debugging
// 	p.API.LogDebug("Plugin configuration", "config", pluginConfig)

// 	var cfg configuration
// 	if err := cfg.FromMap(pluginConfig); err != nil {
// 		return nil, err
// 	}

// 	cfg.Defaults()

// 	updatedPluginConfig, err := cfg.ToMap()
// 	if err != nil {
// 		return nil, err
// 	}
// 	p.API.LogDebug("Updated plugin configuration", "config", updatedPluginConfig)

// 	newCfg.PluginSettings.Plugins["antivirus"] = updatedPluginConfig

// 	return newCfg, nil
// }
