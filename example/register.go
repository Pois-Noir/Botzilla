package main

import (
	"fmt"

	"github.com/Pois-Noir/Botzilla/pkg/component"
)

func main() {

	c1, err := component.NewComponent("c1", "ppap", 4000)

	c1.OnMessage = func(m map[string]string) (map[string]string, error) {
		fmt.Println(m)
		response := map[string]string{}
		response["status"] = "f u"
		return response, nil
	}

	if err != nil {
		fmt.Println(err)
	}
	request := map[string]string{}
	request["Hello"] = "World"

	c2, err := component.NewComponent("c2", "ppap", 4001)
	if err != nil {
		fmt.Println(err)
	}

	c2.OnMessage = func(m map[string]string) (map[string]string, error) {
		fmt.Println(m)
		response := map[string]string{}
		response["status"] = "meow"
		return response, nil
	}

	res, err := c1.SendMessage("c2", request)

	fmt.Println(res)

	for {
	}
}
