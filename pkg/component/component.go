package component

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	global_configs "github.com/Pois-Noir/Botzilla-Utils/global_configs"
	"github.com/Pois-Noir/Botzilla/pkg/core"
)

type Component struct {
	Name       string
	OnMessage  func(map[string]string) (map[string]string, error)
	tunnels    []*tunnel
	serverAddr string
	key        []byte
}

func NewComponent(ServerAddr string, secretKey string, name string, port int) (*Component, error) {

	payload := map[string]string{
		"name": name,
		"port": strconv.Itoa(port),
	}

	encodedPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	// send the message and wait for the response
	_, err = core.Request(
		ServerAddr,
		encodedPayload,
		global_configs.REGISTER_COMPONENT_OPERATION_CODE,
		[]byte(secretKey),
	)

	if err != nil {
		return nil, err
	}

	// generate component with empty message handler
	comp := &Component{
		Name: name,
		OnMessage: func(m map[string]string) (map[string]string, error) {
			return make(map[string]string), nil
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
		go core.ConnectionHandler(conn, key, c.OnMessage)
	}
}

func (c *Component) GetComponents() ([]string, error) {

	rawServerResponse, err := core.Request(
		c.serverAddr,
		[]byte{},
		global_configs.GET_COMPONENTS_OPERATION_CODE,
		c.key,
	)

	if err != nil {
		return nil, err
	}

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

	// send request to server
	rawServerResponse, err := core.Request(
		c.serverAddr,
		destinationBytes,
		global_configs.GET_COMPONENT_OPERATION_CODE,
		c.key,
	)

	if err != nil {
		return nil, err
	}

	// Parsing server response to tcp address
	destinationAddress := string(rawServerResponse)

	// Encoding message content
	encodedBody, err := json.Marshal(message)

	if err != nil {
		return nil, err
	}

	// send request to other component
	rawComponentResponse, err := core.Request(
		destinationAddress,
		encodedBody,
		global_configs.USER_MESSAGE_OPERATION_CODE,
		c.key,
	)

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
