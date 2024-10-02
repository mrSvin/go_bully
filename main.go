package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <node_id>")
		return
	}

	nodeID, err := strconv.Atoi(os.Args[1])
	if err != nil || nodeID < 1 || nodeID > 3 {
		fmt.Println("Node ID must be 1, 2, or 3")
		return
	}

	node := &Node{
		id:   nodeID,
		port: fmt.Sprintf("800%d", nodeID),
	}

	node.start()

	// Блокируем завершение программы
	select {}
}
