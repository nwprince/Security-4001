package main

import (
	"node"
	"time"
)

func main() {
	n := node.New()
	n.WaitTime = 1000 * time.Millisecond
	n.Run("http://localhost:8080")
}
