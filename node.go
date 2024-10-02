package main

import (
	"log"
	"net"
	"time"
)

type Node struct {
	id         int
	port       string
	leaderID   int
	lastLeader time.Time
}

func (n *Node) start() {
	addr, err := net.ResolveUDPAddr("udp", ":"+n.port)
	if err != nil {
		log.Fatalf("Could not resolve address: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Could not listen on port %s: %v", n.port, err)
	}

	go n.listen(conn)

	// Start election process
	n.startElection()

	// Periodically check leader status
	go n.checkLeaderStatus()

	// Periodically send heartbeat if leader
	go n.sendHeartbeat()

	// Блокируем завершение программы, чтобы соединение оставалось открытым
	select {}
}
