package core

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"

	utils_header "github.com/Pois-Noir/Botzilla-Utils/header"
)

func connectionHandler(conn net.Conn, key string, MessageHandler func(map[string]string) (map[string]string, error)) {

	// karim design
	// the decoder needs a conn object
	// i have the conn object as the parameter
	// use the conn object to create a decoder
	// the conn will be converted to a buffered reader by the decoder

	// very good
	defer conn.Close()

	rcvdMsgHeader, err := utils_header.DecodeHeader(conn)

	if err != nil {
		// TODO, header decoding problem
		// we have to send the sender back a error message
		// for retransmission
		// send back error to sender
		// send back and error message
	}

	// check the header

	// getting message lengthg
	requestSize := int32(requestHeader[0]) | // Convert Response Header to int32
		int32(requestHeader[1])<<8 |
		int32(requestHeader[2])<<16 |
		int32(requestHeader[3])<<24

	// reading request
	rawRequest := make([]byte, requestSize)
	// buffered Reader
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
