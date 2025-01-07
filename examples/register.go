package main

import (
	botzilla "botzilla/pkg"
	"fmt"
	"os"
)

type myListener struct{}

func (l myListener) Message(body map[string]string, _ string) (map[string]string, error) {

	fmt.Println("Running command listener")
	fmt.Println(body)
	response := map[string]string{}
	response["darren"] = "gg"
	return response, nil
}

func main() {

	amirgay := myListener{}

	token, err := botzilla.RegisterComponent("localhost:6985", "comp", 6960, amirgay)

	if err != nil {
		fmt.Println("There was an error running example, ", err)
		os.Exit(1)
	}

	println(token)

	response, err := botzilla.GetComponents("localhost:6985", []byte(token))
	fmt.Println(response)

	message := map[string]string{}
	message["umar"] = "is gay"

	res, err := botzilla.SendMessage("localhost:6985", []byte(token), "comp", message)
	fmt.Println(res)

}
