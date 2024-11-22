package core

import (
	"bufio"
	"fmt"
	"net"
)

func BaseCommandHandler(listener func(body string) (string, error)) func(conn net.Conn) {

	return func(conn net.Conn) {

		defer conn.Close()

		// Create a buffered reader
		reader := bufio.NewReader(conn)

		// Read the entire message (this will read until it finds a newline or EOF)
		rawMessage, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading Command:", err)
			}
			return
		}

		// Process the raw message using the listener
		response, err := listener(rawMessage)
		if err != nil {
			fmt.Println("Listener Failed with following err:", err)
			return
		}

		// Send back the response to the client
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error sending response:", err)
		}

	}

}

func BaseMessageHandler(listener func(body string) error) func(conn net.Conn) {

	return func(conn net.Conn) {

		defer conn.Close()

		// Create a buffered reader
		reader := bufio.NewReader(conn)

		// Read the entire message (this will read until it finds a newline or EOF)
		rawMessage, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("Error reading Command:", err)
			}
			return
		}

		err = listener(string(rawMessage))
		if err != nil {
			fmt.Println("Listener Failed with following err:", err)
		}

	}

}

func BaseStreamHandler(listener func()) func(conn net.Conn) {

	return func(conn net.Conn) {
		fmt.Print("Stream handler were the friends we made all along the way")
	}
}
