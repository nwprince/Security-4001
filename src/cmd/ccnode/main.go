package main

import (
	"node"
	"time"
)

func main() {
	n := node.New()
	n.WaitTime = 1000 * time.Millisecond
	n.Run("http://r06nwpkbd.device.mst.edu:8080")
}
