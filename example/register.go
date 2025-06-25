package main

import (
	"fmt"
	botzilla "github.com/Pois-Noir/Botzilla"
)

func main() {

	c1, err := botzilla.NewComponent("c1", "ppap")
	if err != nil {
		fmt.Println(err)
	}

	c1.OnMessage = func(m map[string]string) (map[string]string, error) {
		fmt.Println(m)
		response := map[string]string{}
		response["status"] = "f u"
		return response, nil
	}

	request := map[string]string{}
	request["Hello"] = "World"

	c2, err := botzilla.NewComponent("c2", "ppap")
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
