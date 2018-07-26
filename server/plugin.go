package main

import (
	"io"
	"time"

	"github.com/dutchcoders/go-clamd"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	ClamavHostPort     string
	ScanTimeoutSeconds int
}

func (p *Plugin) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	av := clamd.NewClamd("tcp://" + p.ClamavHostPort)
	abortScan := make(chan bool)
	response, err := av.ScanStream(file, abortScan)
	if err != nil {
		p.API.LogError("Error while scanning for viruses. " + err.Error())
		return nil, "Error while scanning for viruses."
	}
	for {
		select {
		case scanResult, ok := <-response:
			if !ok {
				return info, ""
			}
			if scanResult.Status != clamd.RES_OK {
				p.API.LogWarn("Virus found in file.", "filename", info.Name, "user", info.CreatorId, "scan_result", scanResult.Raw)
				return nil, "Virus found in file."
			}
			continue
		case <-time.After(time.Duration(p.ScanTimeoutSeconds) * time.Second):
			close(abortScan)
			p.API.LogError("Scan timed out.", "filename", info.Name)
			return nil, "Problem with antivirus scanner."
		}
	}
}
