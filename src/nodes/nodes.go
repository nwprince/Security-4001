package nodes

import (
	"encoding/json"
	"fmt"
	"log"
	"messages"
	"os"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Nodes contains all current nodes
var Nodes = make(map[uuid.UUID]*node)

type node struct {
	ID           uuid.UUID
	Platform     string
	Architecture string
	UserName     string
	UserGUID     string
	HostName     string
	Ips          []string
	Pid          int
	Log          *os.File
	FirstTime    time.Time
	MostRecent   time.Time
}

// First handles the first connection that is made
func First(j messages.Base) {
	var sysInfo messages.SysInfo
	sysInfoPayload, err := json.Marshal(j.Payload)

	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(sysInfoPayload, &sysInfo)
	if err != nil {
		log.Panic(err)
	}

	path, _ := os.Getwd()
	dataDir := filepath.Join(path, "data")
	nodeDir := filepath.Join(dataDir, "nodes")

	if _, errD := os.Stat(dataDir); os.IsNotExist(errD) {
		os.Mkdir(dataDir, os.ModeDir)
	}

	if _, errD := os.Stat(nodeDir); os.IsNotExist(errD) {
		os.Mkdir(nodeDir, os.ModeDir)
	}

	nodeDir = filepath.Join(nodeDir, j.ID.String())

	var f *os.File
	if _, err := os.Stat(nodeDir); os.IsNotExist(err) {
		os.Mkdir(nodeDir, os.ModeDir)
		f, err = os.Create(filepath.Join(nodeDir, "log.txt"))
	} else {
		f, err = os.OpenFile(filepath.Join(nodeDir, "log.txt"), os.O_APPEND|os.O_WRONLY, 0600)
	}

	if err != nil {
		log.Panic(err)
	}

	Nodes[j.ID] = &node{
		ID:           j.ID,
		Platform:     sysInfo.Platform,
		Architecture: sysInfo.Architecture,
		UserName:     sysInfo.UserName,
		UserGUID:     sysInfo.UserGUID,
		Ips:          sysInfo.Ips,
		Pid:          sysInfo.Pid,
		Log:          f,
		HostName:     sysInfo.HostName,
		FirstTime:    time.Now(),
		MostRecent:   time.Now(),
	}
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]First conn with %s\r\n", time.Now(), j.ID))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]Platform: %s\r\n", time.Now(), sysInfo.Platform))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]Arch: %s\r\n", time.Now(), sysInfo.Architecture))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]HostName: %s\r\n", time.Now(), sysInfo.HostName))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]Username: %s\r\n", time.Now(), sysInfo.UserName))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]UserGUID: %s\r\n", time.Now(), sysInfo.UserGUID))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]PID: %d\r\n", time.Now(), sysInfo.Pid))
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]IPs: %s\r\n", time.Now(), sysInfo.Ips))
}

// CheckUp will update the log
func CheckUp(j messages.Base) {
	_, ok := Nodes[j.ID]
	if !ok {
		// TODO - do stuff
		return
	}
	Nodes[j.ID].Log.WriteString(fmt.Sprintf("[%s]Node check in\r\n", time.Now()))
	Nodes[j.ID].MostRecent = time.Now()
	return
}
