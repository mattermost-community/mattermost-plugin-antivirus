package main

import (
	"io"
	"sync"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"

	"github.com/IntelXLabs-LLC/go-clamd"
)

type Plugin struct {
	plugin.MattermostPlugin

	configurationLock sync.RWMutex

	configuration *configuration

	// Session tracking
	sessionToConn   map[string]string
	sessionToConnMu sync.RWMutex
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, _ io.Writer) (*model.FileInfo, string) {
	config := p.getConfiguration()

	var av *clamd.Clamd
	if config.ConnectionType == "tcp" {
		av = clamd.NewClamd("tcp://" + config.ClamavHostPort)
	} else {
		av = clamd.NewClamd(config.ClamavSocketPath)
	}
	abortScan := make(chan bool)

	connectionID, found := p.GetConnectionIDForSession(c.SessionId)
	if !found {
		connectionID = ""
		p.API.LogWarn("Session ID not found for user", "session_id", c.SessionId)
	}

	if err := p.API.SendToastMessage(info.CreatorId, connectionID, config.ToastMessageScanning, model.SendToastMessageOptions{
		Position: "bottom-center",
	}); err != nil {
		p.API.LogError("Error while sending toast message. " + err.Error())
	}

	response, err := av.ScanStream(file, abortScan)
	if err != nil {
		p.API.LogError("Error while scanning for viruses. " + err.Error())
		return nil, "File Scanning Server unreachable, contact your Mattermost administrator for assistance."
	}
	for {
		select {
		case scanResult, ok := <-response:
			if !ok {
				if err := p.API.SendToastMessage(info.CreatorId, connectionID, config.ToastMessageSuccess, model.SendToastMessageOptions{
					Position: "bottom-center",
				}); err != nil {
					p.API.LogError("Error while sending success toast message. " + err.Error())
				}
				return info, ""
			}
			if scanResult.Status != clamd.RES_OK {
				p.API.LogWarn("The antivirus service would not allow you to attach this file.", "filename", info.Name, "user", info.CreatorId, "scan_result", scanResult.Raw)
				return nil, "The antivirus service did not allow you to attach this file."
			}
			continue
		case <-time.After(time.Duration(config.ScanTimeoutSeconds) * time.Second):
			close(abortScan)
			p.API.LogError("Scan timed out.", "filename", info.Name)
			return nil, "Problem with antivirus scanner."
		}
	}
}

func (p *Plugin) OnActivate() error {
	p.initializeSessionTracking()
	return nil
}
