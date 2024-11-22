package botzillaclient

import (
	"botzillaclient/core"
	"fmt"
	"net"
)

// Returns a token from server
func StartListener(serverAddress string, config Config, listener Listener) (string, error) {

	// TODO
	/*
		Request Server for a token
	*/
	token := "1234"

	commandHandler := core.BaseCommandHandler(listener.Command)
	messageHandler := core.BaseMessageHandler(listener.Message)
	streamHandler := core.BaseStreamHandler(listener.Stream)

	go core.StartTCPServer(config.CommandPort, commandHandler)
	go core.StartTCPServer(config.MessagePort, messageHandler)
	go core.StartTCPServer(config.StreamPort, streamHandler)

	return token, nil

}

func SendCommand(serverAddress string, token string, follower string, body string) (string, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return "", fmt.Errorf("Error opening connection to botzilla, ", err)
	}

}

func BoardCastMessage(serverAddress string, token string, followers []string, body string) error {

}

func AssignGroup(serverAddress string, token string, groupName string) error {

}

func RemoveGroup(serverAddress string, token string, groupName string) error {

}

func GetAssignedGroups(serverAddress string, token string) ([]string, error) {

}
