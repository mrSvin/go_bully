package main

import (
	"fmt"
	"time"
)

func (n *Node) startElection() {
	fmt.Printf("Node %d: Starting election\n", n.id)
	for i := n.id + 1; i <= 3; i++ {
		n.sendMessage(8000+i, "ELECTION")
	}

	// Wait for OK messages
	time.Sleep(2 * time.Second)

	// If no OK received, declare self as leader
	if n.leaderID == 0 {
		n.leaderID = n.id
		n.lastLeader = time.Now()
		fmt.Printf("Node %d: Declaring self as leader\n", n.id)
		for i := 1; i <= 3; i++ {
			if i != n.id {
				n.sendMessage(8000+i, "LEADER")
			}
		}
	}
}

func (n *Node) checkLeaderStatus() {
	for {
		time.Sleep(5 * time.Second)
		if n.leaderID != n.id && time.Since(n.lastLeader) > 8*time.Second {
			fmt.Printf("Node %d: Leader %d is not responding, starting new election\n", n.id, n.leaderID)
			n.leaderID = 0
			n.startElection()
		}
	}
}
