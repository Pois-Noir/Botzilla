package botzillaclient

import (
	"encoding/json"
	"net"
)

type ports struct {
	command int
	message int
	stream  int
}

func Start(name string, command int, message int, stream int) error {
	return nil
}

func SendCommand(follower string, body string) (string, error) {

	packet := map[string]string{}

	packet["follower"] = follower
	packet["body"] = body

	data, err := json.Marshal(packet)
	if err != nil {
		return "", nil
	}

	conn, err := net.Dial("tcp", "")
	if err != nil {
		return "", nil
	}

	conn.Write(data)

	return "", nil
}

func BroadcastMessage(followers []string, body string) error {
	return nil
}

func StartCommandListener(listener CommandListener, port int) error {
	return nil
}

func StartMessageListener(listener MessageListener, port int) error {
	return nil
}
