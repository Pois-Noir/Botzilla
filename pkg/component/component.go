package component

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	utils_message "github.com/Pois-Noir/Botzilla-Utils/message"
	"github.com/Pois-Noir/Botzilla/pkg/core"
)

type Component struct {
	Name           string
	MessageHandler func(map[string]interface{}) (map[string]interface{}, error)
	tunnels        []*tunnel
	serverAddr     string
	key            []byte
}

func NewComponent(ServerAddr string, secretKey string, name string, port int) (*Component, error) {

	// depricated
	// Generate request content to server
	// create a function for this
	// compSetting := map[string]string{}
	//compSetting["operation_code"] =
	// compSetting["name"] = name
	// compSetting["port"] = strconv.Itoa(port)
	// // include our encoder
	// encodedCompsetting, err := json.Marshal(compSetting)

	// // Operation code 0 is for registration
	// operationCode := []byte{0}
	// // message := append(operationCode, encodedCompsetting...)

	// // send request to server
	// response, err := core.Request(ServerAddr, message, []byte(secretKey))

	compSetting := map[string]interface{}{
		"name": name,
		"port": port,
	}
	// a new message with status code, operation code, and the actual payload
	// TODO get the appropriate operation code or status code
	message := utils_message.NewMessage(0, 0, compSetting)

	// send the message and wait for the response
	response, err := core.Request(ServerAddr, message, []byte(secretKey))

	// need to speak to amir about the response from the server
	// check response from server
	if err != nil {
		return nil, err
	}
	// TODO fix this later
	if string(response) != "registered" {
		return nil, errors.New(string(response))
	}

	// generate component with empty message handler
	comp := &Component{
		Name: name,
		MessageHandler: func(m map[string]interface{}) (map[string]interface{}, error) {
			return make(map[string]interface{}), nil
		},
		key:        []byte(secretKey),
		serverAddr: ServerAddr,
		tunnels:    make([]*tunnel, 0),
	}

	// run tcp listener
	go comp.startListener(port, secretKey)

	return comp, nil

}

func (c *Component) startListener(port int, key string) error {

	// start tcp listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("There was an error starting the server: \n", err)
		return err
	}
	defer listener.Close()

	for {

		// for each connection add a goroutine
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}
		go core.ConnectionHandler(conn, key, c.MessageHandler)
	}
}

func (c *Component) GetComponents() ([]string, error) {

	// // Operation code 2 is for Get All Component
	// operationCode := []byte{69}

	// // send request to server
	// rawServerResponse, err := core.Request(c.serverAddr, operationCode, c.key)

	// create the message to send the server
	message := utils_message.NewMessage(0, 69, make(map[string]interface{}))

	rawServerResponse, err := core.Request(c.serverAddr, message, c.key)
	// parse server response
	var serverResponse []string
	err = json.Unmarshal(rawServerResponse, &serverResponse)
	if err != nil {
		return nil, err
	}

	return serverResponse, nil

}

func (c *Component) SendMessage(componentName string, message map[string]string) (map[string]string, error) {

	// Generate request content to server
	destinationBytes := []byte(componentName)

	// Operation code 2 is for Get Component
	operationCode := []byte{2}
	serverMessage := append(operationCode, destinationBytes...)

	// send request to server
	rawServerResponse, err := core.Request(c.serverAddr, serverMessage, c.key)
	if err != nil {
		return nil, err
	}

	// TODO!!!
	// Server response has to be checked

	// Parsing server response to tcp address
	destinationAddress := string(rawServerResponse)

	// Encoding message content
	encodedBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	// send request to other component
	rawComponentResponse, err := core.Request(destinationAddress, encodedBody, c.key)
	if err != nil {
		return nil, err
	}

	// parse component response
	var componentResponse map[string]string
	err = json.Unmarshal(rawComponentResponse, &componentResponse)
	if err != nil {
		return nil, err
	}

	return componentResponse, nil
}

// TODO
func (c *Component) StartStream(streamName string, input chan byte, port int) error {

	new_tunnel := newTunnel(streamName, input, port)

	c.tunnels = append(c.tunnels, new_tunnel)

	go new_tunnel.start()

	return nil
}

// TODO
func (c *Component) GetComponentStreams(componentName string) error {
	return nil
}

// TODO
func (c *Component) SubscribeStream(componentName string, streamName string) error {
	return nil
}
