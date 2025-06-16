package core

import (
	"fmt"
	"net"

	// for buffered reader
	"bufio"
	"io"

	"github.com/Pois-Noir/Botzilla-Utils/global_configs"
	global_configs "github.com/Pois-Noir/Botzilla-Utils/global_configs"
	header "github.com/Pois-Noir/Botzilla-Utils/header"
	hmac "github.com/Pois-Noir/Botzilla-Utils/hmac"
	"github.com/Pois-Noir/Mammad/decoder"
	"github.com/Pois-Noir/Mammad/encoder"
)

func ConnectionHandler(conn net.Conn, key string, MessageHandler func(map[string]string) (map[string]string, error)) {

	defer conn.Close()

	// creating a buffered Reader
	bReader := bufio.NewReader(conn)
	// read the header from the connection
	// decode it and get a header struct

	var headerBuffer [global_configs.Header_LENGTH]byte
	// TODO
	// CHeck the error
	_, err := io.ReadFull(bReader, headerBuffer[:])
	header, err := header.Decode(bReader)

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
	n, err = io.ReadFull(bReader, hash[:])
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

	convertedRequestBody := convertToStringMap(requestBody)

	if err != nil {
		// tell the server
		// maybe internal server error
	}

	// run the users callback
	responsePayload, err := MessageHandler(convertedRequestBody)
	if err != nil {
		fmt.Println("Error processing message: \n", err)
		return
	}

	encoder = encoder.NewEncoder(responsePayload)
	responsePayloadBuffer := encoder.encode()

	responseHeader := header.NewHeader(
		global_configs.OK_STATUS,
		global_configs.USER_MESSAGE_OPERATION_CODE,
		len(responsePayloadBuffer),
	)

	responseHeaderBuffer, err := responseHeader.Encode()
	if err != nil {
		// speaking with amir
		// call amir 438 282 3324
	}

	response := append(responseHeaderBuffer, responsePayloadBuffer)

	_, err = conn.Write(responseBytes)
	if err != nil {

	}

}

func convertToStringMap(input map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, val := range input {
		result[key] = fmt.Sprint(val)
	}
	return result
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
func Request(serverAddress string, payload []byte, key []byte, operationCode byte) ([]byte, error) {

	// 1. Start TCP connection
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// messageBytes, err := message.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode message: %w", err)
	}

	// 3. Generate HMAC of the entire message (header + payload)
	hash := hmac.GenerateHMAC(payload, key)

	requestHeader := header.NewHeader(global_configs.OK_STATUS, operationCode, len(message))
	encodedHeader := requestHeader.Encode()

	message := append(encodedHeader, payload)

	// TODO
	// Make one io call for speed
	// 4. Send message and HMAC
	_, err = conn.Write(message)
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	_, err = conn.Write(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to write hash: %w", err)
	}

	// 5. Read 4-byte response header
	var responseHeaderBuffer [global_configs.HEADERLENGTH]byte
	_, err = io.ReadFull(conn, responseHeader[:])
	if err != nil {
		return nil, fmt.Errorf("failed to read response header: %w", err)
	}

	responseHeader := header.Decode(responseHeaderBuffer)
	// TODO
	// Check if status was ok

	// 7. Read response body
	rawResponse := make([]byte, responseHeader.Length)
	_, err = io.ReadFull(conn, rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return rawResponse, nil
}
