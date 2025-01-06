package botzilla

import (
	"botzilla/core"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//Hidden Listener/Receiver for server to check if the component is live

// Returns a token from server
func RegisterComponent(serverAddress string, name string, port int, userHandler UserHandler) (string, error) {

	code := []byte{0}
	nameBytes := []byte(name)
	message :=
		append(code, nameBytes...)

	rawResponse, err := requestComponent(serverAddress, message)
	if err != nil {
		return "", err
	}

	token := string(rawResponse)

	//---------------------------------------------------
	//starting listener
	go startListener(port, userHandler)

	return token, nil

}

func startListener(port int, userHandler UserHandler) {

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer listener.Close()

	if err != nil {
		fmt.Println("There was an error starting the server: \n", err)
		return
	}

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: \n", err)
			continue
		}

		go connectionHandler(conn, userHandler)
	}
}

func connectionHandler(conn net.Conn, userHandler UserHandler) {
	defer conn.Close()

	// Create a buffered reader
	reader := bufio.NewReader(conn)

	// Read the entire message (this will read until it finds a newline or EOF)
	rawMessage, err := reader.ReadBytes('\n')

	if err != nil {
		fmt.Printf("Failed to read message: %v\n", err)
		return
	}

	request, err := core.Decode(rawMessage)
	if err != nil {
		fmt.Println("go fuck yourself")
	}

	pt := (*request).Header["type"]

	if pt == "message" {
		response, err := userHandler.Message((*request).Body, (*request).Header["origin"])
		if err != nil {
			fmt.Println("error in user handler")
			fmt.Println(err)
			return
		}
		conn.Write([]byte(response))
	} else if pt == "broadcast" {
		err := userHandler.Broadcast((*request).Body, (*request).Header["origin"])
		if err != nil {
			fmt.Println("error in user handler")
			fmt.Println(err)
			return
		}
	}

}

func SendMessage(serverAddress string, token string, destination string, body map[string]string) (map[string]string, error) {

	code := []byte{2}
	tokenBytes := []byte(token)
	destinationBytes := []byte(destination)

	message :=
		append(
			append(code, tokenBytes...),
			destinationBytes...,
		)

	rawMessage, err := requestServer(serverAddress, message, token)

	destinationAddress := string(rawMessage)

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	rawResponse, err := requestComponent(destinationAddress, encodedBody)
	if err != nil {
		return nil, err
	}

	var decodeMessage map[string]string
	err = json.Unmarshal(rawResponse, &decodeMessage)

	return decodeMessage, nil
}

/*

func BroadCast(serverAddress string, token string, dest []string, body string) error {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}

	message := map[string]string{
		"body":     body,
	}

	decodedMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	conn.Write(decodedMessage)
	conn.Write([]byte("\n"))

	bufferreader := bufio.NewReader(conn)

	rawresponse,err := bufferreader.ReadString('\n')

	if err != nil {
		return err
	}

	return nil
}

func AssignGroup(serverAddress string, token string, groupName string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}

	defer conn.Close()

	message :=
		"0001"+ groupName

	conn.Write([]byte(message))
	conn.Write([]byte("\n"))

	bufferreader := bufio.NewReader(conn)

	res, err := bufferreader.ReadString('\n')

	if err != nil {
		return err
	}

	return nil
}

func RemoveGroup(serverAddress string, token string, groupName string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return err
	}

	defer conn.Close()

	message :=
		"0002"+ groupName

	conn.Write([]byte(message))
	conn.Write([]byte("\n"))

	bufferreader := bufio.NewReader(conn)

	res, err := bufferreader.ReadString('\n')

	if err != nil {
		return err
	}

	return nil
}

func GetAssignedGroups(serverAddress string, token string) ([]string, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil,err
	}

	defer conn.Close()

	message :=
		"0003"

	conn.Write([]byte(message))
	conn.Write([]byte("\n"))

	bufferreader := bufio.NewReader(conn)
	groups:= [] string{}
	groups, err = bufferreader.ReadString('\n')

	if err != nil {
		return nil,err
	}


	return []string{}, nil
}

*/

func GetComponents(serverAddress string, token string) ([]string, error) {

	message := []byte{69}

	rawResponse, err := requestServer(serverAddress, message, token)

	var decodeMessage []string
	err = json.Unmarshal(rawResponse, &decodeMessage)

	if err != nil {
		return nil, err
	}

	return decodeMessage, nil
}

func requestServer(serverAddress string, message []byte, token string) ([]byte, error) {

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

	header := buf.Bytes()

	// Send token for auth
	conn.Write([]byte(token))
	conn.Write(header)
	conn.Write(message)

	responseHeader := [4]byte{}
	conn.Read(responseHeader[:])

	// Convert Response Header to int32
	responseSize := int32(responseHeader[0]) |
		int32(responseHeader[1])<<8 |
		int32(responseHeader[2])<<16 |
		int32(responseHeader[3])<<24

	buffer := make([]byte, responseSize)

	for {
		// Read exactly `fixedLength` bytes
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return nil, err
		}

	}

	return buffer, nil

}

func requestComponent(componentAddress string, message []byte) ([]byte, error) {

	conn, err := net.Dial("tcp", componentAddress)
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

	// Send Header
	conn.Write(buf.Bytes())
	conn.Write(message)

	responseHeader := [4]byte{}
	conn.Read(responseHeader[:])

	// Convert Response Header to int32
	responseSize := int32(responseHeader[0]) |
		int32(responseHeader[1])<<8 |
		int32(responseHeader[2])<<16 |
		int32(responseHeader[3])<<24

	buffer := make([]byte, responseSize)

	for {
		// Read exactly `fixedLength` bytes
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return nil, err
		}

	}

	return buffer, nil

}
