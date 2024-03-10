package network

import (
	"fmt"
	"net"
)

// accepts incoming tcpconnections
// Only masternode should run the TCPserver
func runTCPServer(listenSock net.Listener) {
	for {
		// Accept incoming connections
		AcceptConnection, err := listenSock.Accept()
		if err != nil {
			fmt.Println("Error in connecting to client:", err)
			continue
		}

		go clientHandler(AcceptConnection)

	}
}

func clientHandler(connection net.Conn) {
	//handle incoming data on connection
}
