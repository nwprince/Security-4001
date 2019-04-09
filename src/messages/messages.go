package messages

import (
	uuid "github.com/satori/go.uuid"
)

// Base is the base JSON Object for HTTP Post payloads
type Base struct {
	ID      uuid.UUID   `json:"id"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload,omitempty"`
}

type Transfer struct {
	FileLocation string `json:"dest"`
	FileBlob     string `json:"blob"`
	IsDownload   bool   `json:"download"`
	Job          string `json:"job"`
}

// SysInfo contains json info about the node's system
type SysInfo struct {
	Platform     string   `json:"platform,omitempty"`
	Architecture string   `json:"architecture,omitempty"`
	UserName     string   `json:"username,omitempty"`
	UserGUID     string   `json:"userguid,omitempty"`
	HostName     string   `json:"hostname,omitempty"`
	Pid          int      `json:"pid,omitempty"`
	Ips          []string `json:"ips,omitempty"`
	WaitTime     string   `json:"waittime,omitempty"`
	MaxRetry     int      `json:"maxretry,omitempty"`
}

type CmdPayload struct {
	Command string `json:"executable"`
	Args    string `json:"args"`
	Job     string `json:"job"`
}

type CmdResults struct {
	Job    string `json:"job"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
