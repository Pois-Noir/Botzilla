package botzilla

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
)

type Component struct {
	Name           string
	MessageHandler func(map[string]string) (map[string]string, error)
	tunnels        []*tunnel
	serverAddr     string
	key            []byte
}

func NewComponent(ServerAddr string, secretKey string, name string, port int) (*Component, error) {

	// Generate request content to server
	compSetting := map[string]string{}
	compSetting["name"] = name
	compSetting["port"] = strconv.Itoa(port)
	encodedCompsetting, err := json.Marshal(compSetting)

	// Operation code 0 is for registration
	operationCode := []byte{0}
	message := append(operationCode, encodedCompsetting...)

	// send request to server
	response, err := request(ServerAddr, message, []byte(secretKey))

	// check response from server
	if err != nil {
		return nil, err
	}
	if string(response) != "registered" {
		return nil, errors.New(string(response))
	}

	// generate component with empty message handler
	comp := &Component{
		Name:           name,
		MessageHandler: func(m map[string]string) (map[string]string, error) { return make(map[string]string), nil },
		key:            []byte(secretKey),
		serverAddr:     ServerAddr,
		tunnels:        make([]*tunnel, 0),
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
		go connectionHandler(conn, key, c.MessageHandler)
	}
}

func (c *Component) SendMessage(componentName string, message map[string]string) (map[string]string, error) {

	// Generate request content to server
	destinationBytes := []byte(componentName)

	// Operation code 2 is for Get Component
	operationCode := []byte{2}
	serverMessage := append(operationCode, destinationBytes...)

	// send request to server
	rawServerResponse, err := request(c.serverAddr, serverMessage, c.key)
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
	rawComponentResponse, err := request(destinationAddress, encodedBody, c.key)
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

func (c *Component) GetComponents() ([]string, error) {

	// Operation code 2 is for Get All Component
	operationCode := []byte{69}

	// send request to server
	rawServerResponse, err := request(c.serverAddr, operationCode, c.key)

	// parse server response
	var serverResponse []string
	err = json.Unmarshal(rawServerResponse, &serverResponse)
	if err != nil {
		return nil, err
	}

	return serverResponse, nil

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

func generateHMAC(data []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key) // 32 bytes
	mac.Write(data)
	return mac.Sum(nil)
}

func verifyHMAC(data []byte, key []byte, hash []byte) bool {
	// Generate HMAC for the provided data using the same key
	generatedHMAC := generateHMAC(data, key)

	// Use subtle.ConstantTimeCompare to securely compare the two HMACs
	return subtle.ConstantTimeCompare(generatedHMAC, hash) == 1
}

func connectionHandler(conn net.Conn, key string, MessageHandler func(map[string]string) (map[string]string, error)) {
	defer conn.Close()

	// Reading request header ( indicates request size )
	requestHeader := [4]byte{}
	_, err := conn.Read(requestHeader[:])
	if err != nil {
		fmt.Println("Error reading header: \n", err)
		return
	}
	requestSize := int32(requestHeader[0]) | // Convert Response Header to int32
		int32(requestHeader[1])<<8 |
		int32(requestHeader[2])<<16 |
		int32(requestHeader[3])<<24

	// reading request
	rawRequest := make([]byte, requestSize)
	_, err = conn.Read(rawRequest)

	// reading request hash
	hash := [32]byte{}
	_, err = conn.Read(hash[:])
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	// verifying the hash
	isValid := verifyHMAC(rawRequest, []byte(key), hash[:])
	if !isValid {
		return
	}

	// parsing request
	var request map[string]string
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		fmt.Println("Error unmarshalling message: \n", err)
		return
	}

	// run the users callback
	response, err := MessageHandler(request)
	if err != nil {
		fmt.Println("Error processing message: \n", err)
		return
	}

	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response: \n", err)
		return
	}

	// Generate Header for server
	ResponseLenght := int32(len(encodedResponse))
	RawResponseHeader := new(bytes.Buffer)
	err = binary.Write(RawResponseHeader, binary.LittleEndian, ResponseLenght) // LittleEndian like umar
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return
	}

	ResponseHeader := RawResponseHeader.Bytes()
	_, err = conn.Write(ResponseHeader)
	_, err = conn.Write(encodedResponse)

	if err != nil {
		fmt.Println("Error sending response: \n", err)
	}

}

func request(serverAddress string, message []byte, key []byte) ([]byte, error) {

	// start tcp call
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Generate Header
	messageLength := int32(len(message))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, messageLength) // LittleEndian like umar
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}
	header := buf.Bytes()

	// Generate Hash
	hash := generateHMAC(message, key)

	// Send token for auth
	_, err = conn.Write(header)
	_, err = conn.Write(message)
	_, err = conn.Write(hash)

	// TODO
	// Might need better error handling here
	if err != nil {
		return nil, err
	}

	// Reading Response Header (indicates response size)
	responseHeader := [4]byte{}
	_, err = conn.Read(responseHeader[:])
	if err != nil {
		return nil, err
	}

	// Parsing Header
	responseSize := int32(responseHeader[0]) |
		int32(responseHeader[1])<<8 |
		int32(responseHeader[2])<<16 |
		int32(responseHeader[3])<<24

	// Reading Response
	rawResponse := make([]byte, responseSize)
	_, err = conn.Read(rawResponse)

	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return nil, err
	}

	return rawResponse, nil
}
