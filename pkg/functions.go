package botzillaclient

import (
	"botzillaclient/core"
	"encoding/json"
	"net"
	"strconv"
	"strings"
)

// Returns a token from server
func StartListener(serverAddress string, config Config, listener Listener) (string, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	message := map[string]string{
		"follower": "0000",
		"body":     "0000" + config.Name + "," + "localhost:" + strconv.Itoa(config.CommandPort),
	}

	decodedMessage, err := json.Marshal(message)

	if err != nil {
		return "", err
	}

	conn.Write(decodedMessage)
	conn.Write([]byte("\n"))

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return "", err
	}

	token := string(buffer)

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
		return "", err
	}

	message := map[string]string{
		"follower": follower,
		"body":     body,
	}

	decoded, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	conn.Write(decoded)
	conn.Write([]byte("\n"))

	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return "", err
	}

	response := string(buffer)

	return response, nil
}

func BoardCastMessage(serverAddress string, token string, followers []string, body string) error {
	return nil
}

func AssignGroup(serverAddress string, token string, groupName string) error {
	return nil
}

func RemoveGroup(serverAddress string, token string, groupName string) error {
	return nil
}

func GetAssignedGroups(serverAddress string, token string) ([]string, error) {
	return []string{}, nil
}

func GetComponents(serverAddress string, token string) ([]string, error) {

	response, err := SendCommand(serverAddress, token, "0000", "0001")
	if err != nil {
		return nil, err
	}
	names := strings.Split(response, ",")

	return names, nil

}
