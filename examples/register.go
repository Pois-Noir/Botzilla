package main

import (
	botzillaclient "botzillaclient/pkg"
	"fmt"
	"os"
)

type myListener struct{}

func (l myListener) Message(body string, _ string ) (string, error) {

	fmt.Println("Running command listener")
	fmt.Println(body)

	return "Command Recieved", nil
}

func (l myListener) Broadcast(body string, _ string) error {

	fmt.Println("Running Message Listener")
	fmt.Println(body)

	return nil
}

func main() {

	

	amirgay := myListener{}

	token, err := botzillaclient.RegisterComponent("localhost:6985","comp",6969 , amirgay)

	if err != nil {
		fmt.Println("There was an error running example, ", err)
		os.Exit(1)
	}

	println(token)

}
