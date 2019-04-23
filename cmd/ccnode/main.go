package main

import (
	"time"

	"github.com/nwprince/Security-4001/node"
)

func main() {
	n := node.New()
	n.WaitTime = 1000 * time.Millisecond
	n.Run("http://localhost:8080")
}
