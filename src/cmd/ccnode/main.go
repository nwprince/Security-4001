package main

import (
	"fmt"
	"node"
	"time"
)

func main() {
	n := node.New()
	fmt.Println(n)
	n.WaitTime = 30000 * time.Millisecond
	n.Run("http://localhost:8080")
}
