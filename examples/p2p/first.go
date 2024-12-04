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
	firstName := "first"
	secName := "sec"
	openCommandPortOn := 6787
	openMessagePortOn := 4432
	openStreamPortOn := 1212

	//******************************************************

	config := botzillaclient.Config{
		Name:        firstName,
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

	components, err := botzillaclient.GetComponents(serverAddress, token)

	if err != nil {
		fmt.Println("Error getting the componets")
		os.Exit(1)
	}

	exist := false
	for _, name := range components {
		if name == secName {
			exist = true
			break
		}
	}

	if !exist {
		fmt.Println("Target Component is not registered in botzilla")
		os.Exit(1)
	}

	botzillaclient.SendCommand(serverAddress, token, secName, "fuck you")

}
