package core

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	// for buffered reader
	"bufio"
	"io"

	utils_header "github.com/Pois-Noir/Botzilla-Utils/header"
	hmac "github.com/Pois-Noir/Botzilla-Utils/hmac"
)

func ConnectionHandler(conn net.Conn, key string, MessageHandler func(map[string]string) (map[string]string, error)) {

	// karim design
	// the decoder needs a conn object
	// i have the conn object as the parameter
	// use the conn object to create a decoder
	// the conn will be converted to a buffered reader by the decoder

	// very good
	defer conn.Close()

	// creating a buffered Reader
	bReader := bufio.NewReader(conn)
	// read the header from the connection
	// decode it and get a header struct
	header, err := utils_header.DecodeHeaderBuffered(bReader)

	if err != nil {
		// if there are errors we will send an appropriate response
		// TODO create errors in the error package related to recieving and sending messages
	}
	requestSize := header.Length
	// reading request
	rawRequest := make([]byte, requestSize)

	// we need to read full
	// there is no guarantee it will read upto the rawRequest
	// depricated
	// _, err = conn.Read(rawRequest)

	// n represents the no of bytes read
	n, err := io.ReadFull(bReader, rawRequest[:])
	if n < int(requestSize) {

		// TOOD create errors in the errors package related to rcving bytes
		// send a response back to the sender
		// telling them the message was corrupted
		// log the error
	}
	if err != nil {

	}

	// reading request hash
	hash := [32]byte{}
	n, err = bReader.Read(hash[:])
	if n < 32 {
		// TODO
		// hash was corrupted
	}
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	// verifying the hash
	isValid := hmac.VerifyHMAC(rawRequest, []byte(key), hash[:])
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

func Request(serverAddress string, message []byte, key []byte) ([]byte, error) {

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
	hash := hmac.GenerateHMAC(message, key)

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
