package botzilla

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
)

//Hidden Listener/Receiver for server to check if the component is live

// Returns a token from server
func RegisterComponent(serverAddress string, name string, port int, userHandler UserHandler) (string, error) {

	code := []byte{0}
	var genericToken [16]byte
	compsetting := map[string]string{}
	compsetting["name"] = name
	compsetting["port"] = strconv.Itoa(port)

	encodedCompsetting, err := json.Marshal(compsetting)
	message := append(code, encodedCompsetting...)

	rawResponse, err := requestServer(serverAddress, message, genericToken[:])
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

		go connectionHandler(conn, userHandler)
	}
}

func connectionHandler(conn net.Conn, userHandler UserHandler) {
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

	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	var message map[string]string
	err = json.Unmarshal(buffer, &message)
	if err != nil {
		fmt.Println("Error unmarshalling message: \n", err)
		return
	}

	response, err := userHandler.Message(message, conn.RemoteAddr().String())
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
	conn.Write(headerHeader)
	_, err = conn.Write(encodedResponse)

}

func SendMessage(serverAddress string, token []byte, destination string, body map[string]string) (map[string]string, error) {

	code := []byte{2}
	destinationBytes := []byte(destination)

	message := append(code, destinationBytes...)

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

func GetComponents(serverAddress string, token []byte) ([]string, error) {

	message := []byte{69}

	rawResponse, err := requestServer(serverAddress, message, token)

	var decodeMessage []string
	err = json.Unmarshal(rawResponse, &decodeMessage)

	if err != nil {
		return nil, err
	}

	return decodeMessage, nil
}

func requestServer(serverAddress string, message []byte, token []byte) ([]byte, error) {

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
	conn.Write(token)
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

	_, err = conn.Read(buffer)

	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return nil, err
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

	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return nil, err
	}

	return buffer, nil

}
