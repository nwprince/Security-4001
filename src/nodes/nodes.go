package nodes

import (
	"encoding/json"
	"errors"
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
	log          *os.File
	channel      chan []Job
	FirstTime    time.Time
	MostRecent   time.Time
	WaitTime     string
	MaxRetry     int
}

type Job struct {
	ID      string
	Type    string
	Status  string
	Args    []string
	Created time.Time
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
		log:          f,
		channel:      make(chan []Job, 10),
		HostName:     sysInfo.HostName,
		FirstTime:    time.Now(),
		MostRecent:   time.Now(),
		WaitTime:     sysInfo.WaitTime,
		MaxRetry:     sysInfo.MaxRetry,
	}
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]First conn with %s\r\n", time.Now(), j.ID))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]Platform: %s\r\n", time.Now(), sysInfo.Platform))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]Arch: %s\r\n", time.Now(), sysInfo.Architecture))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]HostName: %s\r\n", time.Now(), sysInfo.HostName))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]Username: %s\r\n", time.Now(), sysInfo.UserName))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]UserGUID: %s\r\n", time.Now(), sysInfo.UserGUID))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]PID: %d\r\n", time.Now(), sysInfo.Pid))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]IPs: %s\r\n", time.Now(), sysInfo.Ips))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]WaitTime: %s\r\n", time.Now(), sysInfo.WaitTime))
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]MaxRetry: %d\r\n", time.Now(), sysInfo.MaxRetry))
}

// CheckUp will update the log
func CheckUp(j messages.Base) (messages.Base, error) {
	_, ok := Nodes[j.ID]
	if !ok {
		// TODO - do stuff
		log.Panic("help")
	}
	Nodes[j.ID].log.WriteString(fmt.Sprintf("[%s]Node check in\r\n", time.Now()))
	Nodes[j.ID].MostRecent = time.Now()

	if len(Nodes[j.ID].channel) >= 1 {
		job := <-Nodes[j.ID].channel

		m, err := GetMessageForJob(j.ID, job[0])
		return m, err
	}
	m := messages.Base{
		ID:   j.ID,
		Type: "ServerOk",
	}
	return m, nil
}

// GetStatus will return the status of the node
func GetStatus(id uuid.UUID) string {
	var status string
	dur, err := time.ParseDuration(Nodes[id].WaitTime)
	if err != nil {
		log.Println("warn", fmt.Sprintf("Error with conv %s to a time duration: %s", Nodes[id].WaitTime, err.Error()))
	}

	if Nodes[id].MostRecent.Add(dur).After(time.Now()) {
		status = "Active"
	} else if Nodes[id].MostRecent.Add(dur * time.Duration(Nodes[id].MaxRetry+1)).After(time.Now()) {
		status = "Delayed"
	} else {
		status = "Dead"
	}
	return status
}

func AddJob(nodeID uuid.UUID, jobType string, jobArgs []string) (string, error) {
	isNode := false
	for k := range Nodes {
		if Nodes[k].ID == nodeID {
			isNode = true
		}
	}

	if isNode {
		job := Job{
			Type:    jobType,
			Status:  "created",
			Args:    jobArgs,
			Created: time.Now(),
		}

		if nodeID.String() == "ffffffff-ffff-ffff-ffff-ffffffffffff" {

			for k := range Nodes {
				s := Nodes[k].channel
				id, _ := uuid.NewV4()
				job.ID = id.String()
				s <- []Job{job}
				Nodes[k].log.WriteString(fmt.Sprintf("[%s]Created job Type:%s, ID:%s, Status:%s, "+"Args: %s \r\n", time.Now(), job.Type, job.ID, job.Status, job.Args))
			}
			return job.ID, nil
		}

		id, _ := uuid.NewV4()
		job.ID = id.String()
		s := Nodes[nodeID].channel
		s <- []Job{job}
		Nodes[nodeID].log.WriteString(fmt.Sprintf("[%s]Created job Type:%s, ID:%s, Status:%s, "+"Args: %s \r\n", time.Now(), job.Type, job.ID, job.Status, job.Args))
		return job.ID, nil
	}
	return "", errors.New("bad id")
}

func GetMessageForJob(nodeID uuid.UUID, job Job) (messages.Base, error) {
	m := messages.Base{
		ID: nodeID,
	}
	switch job.Type {
	case "cmd":
		m.Type = "CmdPayload"
		p := messages.CmdPayload{
			Command: job.Args[0],
			Job:     job.ID,
		}

		k := marshalMessage(p)
		m.Payload = (*json.RawMessage)(&k)

	}

	return m, nil
}

func marshalMessage(m interface{}) []byte {
	k, _ := json.Marshal(m)
	return k
}
