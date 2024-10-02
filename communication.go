package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func (n *Node) listen(conn *net.UDPConn) {
	defer conn.Close() // Закрываем соединение только после завершения работы listen

	buf := make([]byte, 1024)
	for {
		l, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		message := strings.TrimSpace(string(buf[:l]))
		parts := strings.Split(message, ":")
		if len(parts) < 2 {
			continue
		}

		senderID, _ := strconv.Atoi(parts[0])
		command := parts[1]

		switch command {
		case "ELECTION":
			if senderID < n.id {
				n.sendMessage(addr.Port, "OK")
				n.startElection()
			}
		case "OK":
			// Do nothing, just acknowledge
		case "LEADER":
			n.leaderID = senderID
			n.lastLeader = time.Now()
			fmt.Printf("Node %d: New leader is %d\n", n.id, n.leaderID)
		case "HEARTBEAT":
			if senderID == n.leaderID {
				n.lastLeader = time.Now()
			}
		}
	}
}

func (n *Node) sendMessage(port int, message string) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Printf("Could not resolve address: %v", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("Could not dial UDP: %v", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte(fmt.Sprintf("%d:%s", n.id, message)))
	if err != nil {
		log.Printf("Could not send message: %v", err)
	}
}

func (n *Node) sendHeartbeat() {
	for {
		time.Sleep(3 * time.Second)
		if n.leaderID == n.id {
			for i := 1; i <= 3; i++ {
				if i != n.id {
					n.sendMessage(8000+i, "HEARTBEAT")
				}
			}
		}
	}
}
