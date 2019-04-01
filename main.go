package main

import (
	"cli"
	"encoding/json"
	"fmt"
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
			// TODO: Do something with status
		}

	} else {
		w.WriteHeader(404)
	}

}
