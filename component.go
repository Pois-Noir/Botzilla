package botzilla

import (
	"encoding/json"
	"fmt"
	"github.com/Pois-Noir/Botzilla-Utils/global_configs"
	"net"

	"github.com/grandcat/zeroconf"
)

type Component struct {
	Name      string
	server    *zeroconf.Server
	OnMessage func(map[string]string) (map[string]string, error)
	tunnels   []*tunnel
	key       []byte
}

func NewComponent(name string, secretKey string) (*Component, error) {
	// generate component with empty message handler
	comp := &Component{
		Name: name,
		OnMessage: func(m map[string]string) (map[string]string, error) {
			return make(map[string]string), nil
		},
		key:     []byte(secretKey),
		tunnels: make([]*tunnel, 0),
	}

	// run tcp listener
	port, err := comp.startListening()
	if err != nil {
		return nil, err
	}

	server, err := zeroconf.Register(
		name,
		"_botzilla._tcp",
		"local.",
		port,
		[]string{"id=botzilla_" + name},
		nil,
	)

	if err != nil {
		return nil, err
	}

	comp.server = server

	return comp, nil

}

func (c *Component) startListening() (int, error) {

	// start tcp listener
	listener, err := net.Listen("tcp", ":0") // open on the random free port
	if err != nil {
		return -1, err
	}

	// get the port
	port := listener.Addr().(*net.TCPAddr).Port

	go c.handleTCP(listener)

	return port, nil
}

func (c *Component) handleTCP(listener net.Listener) {
	defer listener.Close()

	for {

		// for each connection add a goroutine
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}
		go ConnectionHandler(conn, c.key, c.OnMessage)
	}
}

func (c *Component) SendMessage(componentName string, message map[string]string) (map[string]string, error) {

	compIP, err := GetComponent(componentName)
	if err != nil {
		return nil, err
	}

	//// Encoding message content
	encodedBody, err := json.Marshal(message)
	//
	if err != nil {
		return nil, err
	}
	//
	//// send request to other component
	rawComponentResponse, err := Request(
		compIP,
		encodedBody,
		global_configs.USER_MESSAGE_OPERATION_CODE,
		c.key,
	)
	//
	if err != nil {
		return nil, err
	}
	//
	//// parse component response
	var componentResponse map[string]string
	err = json.Unmarshal(rawComponentResponse, &componentResponse)
	//
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
