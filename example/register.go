package main

import (
	"fmt"

	"github.com/Pois-Noir/Botzilla/pkg/component"
)

func main() {

	c1, err := component.NewComponent("localhost:6985", "ppap", "c1", 4000)

	c1.OnMessage = func(m map[string]string) (map[string]string, error) {
		fmt.Println(m)
		response := map[string]string{}
		response["status"] = "f u"
		return response, nil
	}

	if err != nil {
		fmt.Println(err)
	}

	c2, err := component.NewComponent("localhost:6985", "ppap", "c2", 4001)

	if err != nil {
		fmt.Println(err)
		fmt.Println("ahhhh")
	}

	request := map[string]string{}
	request["Hello"] = "World"

	response, err := c2.SendMessage("c1", request)

	fmt.Println(response)
}
