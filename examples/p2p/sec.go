package main

import (
	botzillaclient "botzillaclient/pkg"
	"fmt"
	"os"
)

type myListener struct{}

func (l myListener) Command(body string) (string, error) {

	fmt.Println("Running command listener")
	fmt.Println(body)

	return "Command Recieved", nil
}

func (l myListener) Message(body string) error {

	fmt.Println("Running Message Listener")
	fmt.Println(body)

	return nil
}

func (l myListener) Stream() {
	fmt.Println("streaming is not supported yet :O")
}

func main() {

	//******************************************************
	//	SENARIO Setting
	serverAddress := "localhost:6985"
	secName := "sec"
	openCommandPortOn := 6788
	openMessagePortOn := 4433
	openStreamPortOn := 1213

	//******************************************************

	config := botzillaclient.Config{
		Name:        secName,
		CommandPort: openCommandPortOn,
		MessagePort: openMessagePortOn,
		StreamPort:  openStreamPortOn,
	}

	listener := myListener{}

	token, err := botzillaclient.StartListener(serverAddress, config, listener)

	if err != nil {
		fmt.Println("There was an error running example, ", err)
		os.Exit(1)
	}

	println(token)

	for {
	}
}
