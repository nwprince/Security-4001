package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nwprince/Security-4001/cli"
	"github.com/nwprince/Security-4001/messages"
	"github.com/nwprince/Security-4001/nodes"
)

func init() {
	go cli.Shell()
}

func main() {
	http.HandleFunc("/", cellHandler)
	http.HandleFunc("/aarch64", aarch64)
	http.HandleFunc("/armv7a", armv7a)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func aarch64(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, "data/temp/v0.0.1/0.0.1/ccNode-aarch64")
}
func armv7a(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, "data/temp/v0.0.1/0.0.1/ccNode-armv7a")
}

func cellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "POST" {
		var payload json.RawMessage
		j := messages.Base{
			Payload: &payload,
		}
		json.NewDecoder(r.Body).Decode(&j)

		switch j.Type {

		case "first":
			// TODO: Checkin new nodes
			nodes.First(j)
			break
		case "status":
			w.Header().Set("Content-Type", "application/json")
			// TODO: Do something with status
			x, err := nodes.CheckUp(j)
			if err != nil {
				log.Println(err)
			}
			json.NewEncoder(w).Encode(x)

		case "CmdResults":
			var p messages.CmdResults
			json.Unmarshal(payload, &p)
			fmt.Println("")
			if len(p.Stdout) > 0 {
				fmt.Println(j.ID.String() + " results: ~> SUCCESS\r\n" + p.Stdout)
			}
			if len(p.Stderr) > 0 {
				fmt.Println(j.ID.String() + " results: ~> ERROR\r\n" + p.Stderr)
			}

		case "TransferResults":
			var p messages.Transfer
			json.Unmarshal(payload, &p)
			if p.IsDownload {
				str, _ := base64.StdEncoding.DecodeString(p.FileBlob)
				filename := strings.Split(p.FileLocation, "/")
				dir, _ := os.Getwd()
				path := filepath.Join(dir, "data", "nodes", j.ID.String(), filename[len(filename)-1])
				err := ioutil.WriteFile(path, str, 0644)
				if err != nil {
					panic(err)
				}
			}

		}

	} else {
		w.WriteHeader(404)
	}

}
