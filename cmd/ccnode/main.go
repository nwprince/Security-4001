package main

import (
	"fmt"
	"node"
)

func main() {
	n := node.New()
	fmt.Println(n)
	n.Run("http://localhost:8080")
}
