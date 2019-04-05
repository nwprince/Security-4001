package main

import (
	"cli"
	"encoding/json"
	"fmt"
	"log"
	"messages"
	"net/http"
	"nodes"
)

func init() {
	go cli.Shell()
}

func main() {
	http.HandleFunc("/", cellHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func cellHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("endpoint hit")

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
			if len(p.Stdout) > 0 {
				fmt.Println(j.ID.String() + " results: ~> SUCCESS\r\n" + p.Stdout)
			}
			if len(p.Stderr) > 0 {
				fmt.Println(j.ID.String() + " results: ~> SUCCESS\r\n" + p.Stderr)
			}
		}

	} else {
		w.WriteHeader(404)
	}

}
