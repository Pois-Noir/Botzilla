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

	config := botzillaclient.Config{
		Name:        "Comp1",
		CommandPort: 6787,
		MessagePort: 4432,
		StreamPort:  1212,
	}

	listener := myListener{}

	token, err := botzillaclient.StartListener("localhost:6985", config, listener)

	if err != nil {
		fmt.Println("There was an error running example, ", err)
		os.Exit(1)
	}

	println(token)

}
