package client

import (
	"encoding/json"
	"net"
)

type server struct {
	Address string
	Command int
	Message int
	Stream  int
}

type Client struct {
	name string
	port int
	s    *server
}

func NewClient(name string, port int) *Client {
	return &Client{
		name: name,
		port: port,
		s:    nil,
	}
}

func (c *Client) RegisterComponent(address string, serverPort int) {

	// TODO
	/*
		Request server and ask for its ports
	*/
	mPort := 1
	sPort := 2

	(*c).s = &server{
		Address: address,
		Command: serverPort,
		Message: mPort,
		Stream:  sPort,
	}

}

func (c *Client) SendCommand(follower string, body string) (string, error) {

	// TODO
	/*
		Check if server is nil
	*/

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

func (c *Client) BroadcastMessage(followers []string, body string) error {

	// TODO
	/*
		Check if server is nil
	*/

	// TODO
	// Implement it

	return nil
}

func (c *Client) StartCommandListener(listener CommandListener, port int) error {

	// TODO
	/*
		Check if server is nil
	*/

	// TODO
	// Implement it

	return nil
}

func (c *Client) StartMessageListener(listener MessageListener, port int) error {

	// TODO
	/*
		Check if server is nil
	*/

	// TODO
	// Implement it

	return nil
}
