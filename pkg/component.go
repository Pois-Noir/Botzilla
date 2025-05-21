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
	Name       string
	serverAddr string
	key        []byte
}

func NewComponent(ServerAddr string, secretKey string, name string, port int, MessageListener func(map[string]string) (map[string]string, error)) (*Component, error) {

	code := []byte{0}
	compSetting := map[string]string{}
	compSetting["name"] = name
	compSetting["port"] = strconv.Itoa(port)

	encodedCompsetting, err := json.Marshal(compSetting)
	message := append(code, encodedCompsetting...)

	response, err := request(ServerAddr, message, []byte(secretKey))
	if err != nil {
		return nil, err
	}

	if string(response) != "registered" {
		return nil, errors.New(string(response))
	}

	go startListener(port, secretKey, MessageListener)

	return &Component{Name: name, key: []byte(secretKey), serverAddr: ServerAddr}, nil

}

func (c *Component) SendMessage(name string, message map[string]string) (map[string]string, error) {

	code := []byte{2}
	destinationBytes := []byte(name)

	serverMessage := append(code, destinationBytes...)

	rawMessage, err := request(c.serverAddr, serverMessage, c.key)
	if err != nil {
		return nil, err
	}

	destinationAddress := string(rawMessage)

	encodedBody, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	rawResponse, err := request(destinationAddress, encodedBody, c.key)
	if err != nil {
		return nil, err
	}

	var decodeMessage map[string]string
	err = json.Unmarshal(rawResponse, &decodeMessage)
	if err != nil {
		return nil, err
	}

	return decodeMessage, nil
}

func (c *Component) GetComponents() ([]string, error) {

	message := []byte{69}

	rawResponse, err := request(c.serverAddr, message, c.key)

	var decodeMessage []string
	err = json.Unmarshal(rawResponse, &decodeMessage)

	if err != nil {
		return nil, err
	}

	return decodeMessage, nil

}

func (c *Component) Broadcast(message map[string]string) error {

	arr, err := c.GetComponents()
	if err != nil {
		return err
	}

	for _, comp := range arr {
		go c.SendMessage(comp, message)
	}

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

func startListener(port int, key string, userHandler func(map[string]string) (map[string]string, error)) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		fmt.Println("There was an error starting the server: \n", err)
		return
	}

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}

		go connectionHandler(conn, key, userHandler)
	}
}

func connectionHandler(conn net.Conn, key string, userHandler func(map[string]string) (map[string]string, error)) {
	defer conn.Close()

	requestHeader := [4]byte{}

	_, err := conn.Read(requestHeader[:])
	if err != nil {
		fmt.Println("Error reading header: \n", err)
		return
	}

	// Convert Response Header to int32
	requestSize := int32(requestHeader[0]) |
		int32(requestHeader[1])<<8 |
		int32(requestHeader[2])<<16 |
		int32(requestHeader[3])<<24

	buffer := make([]byte, requestSize)

	_, err = conn.Read(buffer)

	hash := [32]byte{}
	_, err = conn.Read(hash[:])

	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	isValid := verifyHMAC(buffer, []byte(key), hash[:])
	if !isValid {
		return
	}

	var message map[string]string
	err = json.Unmarshal(buffer, &message)
	if err != nil {
		fmt.Println("Error unmarshalling message: \n", err)
		return
	}

	response, err := userHandler(message)
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
	messageLength := int32(len(encodedResponse))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, messageLength) // LittleEndian like umar
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return
	}

	headerHeader := buf.Bytes()
	_, err = conn.Write(headerHeader)
	_, err = conn.Write(encodedResponse)

	if err != nil {
		fmt.Println("Error sending response: \n", err)
	}

}

func request(serverAddress string, message []byte, key []byte) ([]byte, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	// Generate Header for server
	messageLength := int32(len(message))
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, messageLength) // LittleEndian like umar
	if err != nil {
		fmt.Println("binary.Write failed:", err)
		return nil, err
	}

	hash := generateHMAC(message, key)

	header := buf.Bytes()

	// Send token for auth
	_, err = conn.Write(header)
	_, err = conn.Write(message)
	_, err = conn.Write(hash)

	if err != nil {
		return nil, err
	}

	responseHeader := [4]byte{}
	_, err = conn.Read(responseHeader[:])
	if err != nil {
		return nil, err
	}

	// Convert Response Header to int32
	responseSize := int32(responseHeader[0]) |
		int32(responseHeader[1])<<8 |
		int32(responseHeader[2])<<16 |
		int32(responseHeader[3])<<24

	buffer := make([]byte, responseSize)

	_, err = conn.Read(buffer)

	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return nil, err
	}

	return buffer, nil
}
