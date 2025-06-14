package core

import (
	"encoding/binary"
	"fmt"
	"net"

	// for buffered reader
	"bufio"
	"io"

	utils_header "github.com/Pois-Noir/Botzilla-Utils/header"
	hmac "github.com/Pois-Noir/Botzilla-Utils/hmac"
	utils_message "github.com/Pois-Noir/Botzilla-Utils/message"
	"github.com/Pois-Noir/Mammad/decoder"
)

func ConnectionHandler(conn net.Conn, key string, MessageHandler func(map[string]interface{}) (map[string]interface{}, error)) {

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
	// get the message length
	requestSize := header.Length
	// buffer to store the message
	rawRequest := make([]byte, requestSize)

	// We use io.ReadFull to guarantee that we read exactly `requestSize` bytes.
	// A normal conn.Read may return fewer bytes than requested if the kernel buffer doesn't contain all data yet.
	// io.ReadFull retries reads internally until the buffer is filled or an error occurs.
	//
	// However, if there's a transmission error (e.g. connection dropped mid-transfer),
	// io.ReadFull may return with an error *and* partial data (n < requestSize).
	// That's why we check `n` to verify how many bytes were actually read — not just rely on the error.
	n, err := io.ReadFull(bReader, rawRequest[:])
	if n < int(requestSize) {

		// TOOD create errors in the errors package related to rcving bytes
		// send a response back to the sender
		// telling them the message was corrupted
		// log the error
	}
	if err != nil {
		// do something
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
	decoder := decoder.NewDecoderBytes(rawRequest)
	requestBody, err := decoder.Decode(len(rawRequest))
	if err != nil {
		// tell the server
		// maybe internal server error
	}

	// run the users callback
	responseMap, err := MessageHandler(requestBody)
	if err != nil {
		fmt.Println("Error processing message: \n", err)
		return
	}

	response := utils_message.NewMessage(0, 0, responseMap)

	responseBytes, err := response.Encode()
	if err != nil {
		// speaking with amir
		// call amir 438 282 3324
	}

	_, err = conn.Write(responseBytes)
	if err != nil {

	}

	// // Generate Header for server
	// ResponseLenght := int32(len(encodedResponse))
	// RawResponseHeader := new(bytes.Buffer)
	// err = binary.Write(RawResponseHeader, binary.LittleEndian, ResponseLenght) // LittleEndian like umar
	// if err != nil {
	// 	fmt.Println("binary.Write failed:", err)
	// 	return
	// }

	// ResponseHeader := RawResponseHeader.Bytes()
	// _, err = conn.Write(ResponseHeader)
	// _, err = conn.Write(encodedResponse)

	// if err != nil {
	// 	fmt.Println("Error sending response: \n", err)
	// }

}

// func Request(serverAddress string, message []byte, key []byte) ([]byte, error) {

// 	// start tcp call
// 	conn, err := net.Dial("tcp", serverAddress)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer conn.Close()

// 	// Generate Header
// 	messageLength := int32(len(message))
// 	buf := new(bytes.Buffer)
// 	err = binary.Write(buf, binary.LittleEndian, messageLength) // LittleEndian like umar
// 	if err != nil {
// 		fmt.Println("binary.Write failed:", err)
// 		return nil, err
// 	}
// 	header := buf.Bytes()

// 	// Generate Hash
// 	hash := hmac.GenerateHMAC(message, key)

// 	// Send token for auth
// 	_, err = conn.Write(header)
// 	_, err = conn.Write(message)
// 	_, err = conn.Write(hash)

// 	// TODO
// 	// Might need better error handling here
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Reading Response Header (indicates response size)
// 	responseHeader := [4]byte{}
// 	_, err = conn.Read(responseHeader[:])
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parsing Header
// 	responseSize := int32(responseHeader[0]) |
// 		int32(responseHeader[1])<<8 |
// 		int32(responseHeader[2])<<16 |
// 		int32(responseHeader[3])<<24

// 	// Reading Response
// 	rawResponse := make([]byte, responseSize)
// 	_, err = conn.Read(rawResponse)

// 	if err != nil {
// 		fmt.Printf("Error reading from connection: %v\n", err)
// 		return nil, err
// 	}

//		return rawResponse, nil
//	}
func Request(serverAddress string, message *utils_message.Message, key []byte) ([]byte, error) {
	// 1. Start TCP connection
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	messageBytes, err := message.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	// 3. Generate HMAC of the entire message (header + payload)
	hash := hmac.GenerateHMAC(messageBytes, key)

	// 4. Send message and HMAC
	_, err = conn.Write(messageBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	_, err = conn.Write(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to write hash: %w", err)
	}

	// 5. Read 4-byte response header (LittleEndian)
	var responseHeader [4]byte
	_, err = io.ReadFull(conn, responseHeader[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read response header: %w", err)
	}

	// 6. Parse response length
	responseSize := binary.LittleEndian.Uint32(responseHeader[:])

	// 7. Read response body
	rawResponse := make([]byte, responseSize)
	_, err = io.ReadFull(conn, rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return rawResponse, nil
}
