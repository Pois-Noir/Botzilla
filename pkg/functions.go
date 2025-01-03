package botzillaclient

import (
	"botzillaclient/core"
	"bufio"
	"fmt"
	"net"
)

//Hidden Listener/Receiver for server to check if the component is live

// Returns a token from server
func RegisterComponent(serverAddress string, name string, port int, userHandler UserHandler) (string, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	//-------------------------------------------------
	// Registration
	message :=
		"0000"+ name
	
	
	conn.Write([]byte(message))
	conn.Write([]byte("\n"))

	// Todo: Buffer size problem for receiving data
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		return "", err
	}

	token := string(buffer)

	//---------------------------------------------------
	//starting listener
	go startListener(port, userHandler)

	return token, nil

}

func startListener(port int, userHandler UserHandler){

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

func connectionHandler(conn net.Conn, userHandler UserHandler){
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
		response, err := userHandler.Message((*request).Body , (*request).Header["origin"])
		if err != nil {
			fmt.Println("error in user handler")
			fmt.Println(err)
			return
		}
		conn.Write([]byte(response))
	} else if pt == "broadcast" {
		err := userHandler.Broadcast((*request).Body , (*request).Header["origin"])
		if err != nil {
			fmt.Println("error in user handler")
			fmt.Println(err)
			return
		}
	}

}

/*
func SendMessage(serverAddress string, token string, dest string, body string) (string, error) {

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return "", err
	}

	//TODO get the address of the destination, then send the message
	message := map[string]string{
		"body":     body,
	}

	decodedMessage, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	conn.Write(decodedMessage)
	conn.Write([]byte("\n"))

	bufferreader := bufio.NewReader(conn)

	rawresponse,err := bufferreader.ReadString('\n')

	if err != nil {
		return "", err
	}

	response := string(rawresponse)

	return response, nil
}

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

func GetComponents(serverAddress string, token string) ([]string, error) {

	response, err := SendCommand(serverAddress, token, "0000", "0001")
	if err != nil {
		return nil, err
	}
	names := strings.Split(response, ",")

	return names, nil

}
*/