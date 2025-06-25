package botzilla

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"

	global_configs "github.com/Pois-Noir/Botzilla-Utils/global_configs"
	header_pkg "github.com/Pois-Noir/Botzilla-Utils/header"
	hmac "github.com/Pois-Noir/Botzilla-Utils/hmac"
)

func ConnectionHandler(conn net.Conn, key []byte, MessageHandler func(map[string]string) (map[string]string, error)) {

	defer conn.Close()

	// creating a buffered Reader
	bReader := bufio.NewReader(conn)

	var headerBuffer [global_configs.HEADER_LENGTH]byte
	n, err := io.ReadFull(bReader, headerBuffer[:])

	// TODO
	if n != global_configs.HEADER_LENGTH {
		return
	}

	// TODO
	if err != nil {
		return
	}

	header, err := header_pkg.Decode(headerBuffer[:])

	// TODO
	if err != nil {
		return
	}
	// get the message length
	requestSize := header.PayloadLength
	RequestPayloadBuffer := make([]byte, requestSize)

	// We use io.ReadFull to guarantee that we read exactly `requestSize` bytes.
	n, err = io.ReadFull(bReader, RequestPayloadBuffer)

	// TODO
	if uint32(n) < requestSize {
		return
	}
	// TODO
	if err != nil {
		return
	}

	// reading request hash
	hash := [32]byte{}
	n, err = io.ReadFull(bReader, hash[:])

	// TODO
	if err != nil {
		return
	}

	// TODO
	if n < 32 {
		// hash was corrupted
	}

	// verifying the hash
	isValid := hmac.VerifyHMAC(RequestPayloadBuffer, key, hash[:])
	// TODO
	if !isValid {
		return
	}

	// parsing request
	RequestPayload := map[string]string{}
	err = json.Unmarshal(RequestPayloadBuffer[:], &RequestPayload)

	// TODO
	if err != nil {
		return
	}

	// run user callback
	ResponsePayload, err := MessageHandler(RequestPayload)

	// TODO
	if err != nil {
		return
	}

	ResponsePayloadBuffer, err := json.Marshal(ResponsePayload)

	// TODO
	if err != nil {
		fmt.Println("Error processing message: \n", err)
		return
	}

	ResponseHeader := header_pkg.NewHeader(
		global_configs.OK_STATUS,
		global_configs.USER_MESSAGE_OPERATION_CODE,
		uint32(len(ResponsePayloadBuffer)),
	)
	ResponseHeaderBuffer := ResponseHeader.Encode()

	ResponseBuffer := append(ResponseHeaderBuffer, ResponsePayloadBuffer...)
	_, err = conn.Write(ResponseBuffer)

	if err != nil {
		return
	}

}

func Request(serverAddress string, payload []byte, operationCode byte, key []byte) ([]byte, error) {

	// 1. Start TCP connection
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 3. Generate HMAC of the entire message (header + payload)
	hash := hmac.GenerateHMAC(payload, key)

	requestHeader := header_pkg.NewHeader(
		global_configs.OK_STATUS,
		operationCode,
		uint32(len(payload)),
	)
	encodedHeader := requestHeader.Encode()

	message := append(encodedHeader, payload...)

	// Make one io call for speed
	// 4. Send message and HMAC
	_, err = conn.Write(message)

	// TODO
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	_, err = conn.Write(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to write hash: %w", err)
	}

	// 5. Read 4-byte response header
	var responseHeaderBuffer [global_configs.HEADER_LENGTH]byte
	_, err = io.ReadFull(conn, responseHeaderBuffer[:])
	// TODO
	if err != nil {
		return nil, fmt.Errorf("failed to read response header: %w", err)
	}

	responseHeader, err := header_pkg.Decode(responseHeaderBuffer[:])
	// TODO
	if err != nil {
		return nil, err
	}
	// Check if status was ok

	// 7. Read response body
	rawResponse := make([]byte, responseHeader.PayloadLength)
	_, err = io.ReadFull(conn, rawResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return rawResponse, nil
}
