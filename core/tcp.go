package core

import (
	"fmt"
	"net"
	"os"
)

func StartTCPServer(port int, handler func(conn net.Conn)) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("There was an error starting the server: \n", err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Commanding server has started on port 8080")

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}

		go handler(conn)
	}
}
