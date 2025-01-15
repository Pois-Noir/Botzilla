package main

import (
	botzilla "botzilla/pkg"
	"fmt"
)

func Message(body map[string]string) (map[string]string, error) {

	fmt.Println("Running command listener")
	fmt.Println(body)
	response := map[string]string{}
	response["darren"] = "gg"
	return response, nil
}

func main() {

	compA, err := botzilla.NewComponent("localhost:6985", "ppap", "comp1", 6960, Message)
	if err != nil {
	}

	compB, err := botzilla.NewComponent("localhost:6985", "ppap", "comp2", 6942, Message)

	fmt.Println(compA.GetComponents())

	m := map[string]string{}
	m["data"] = "a request from comp1 to comp2"
	response, err := compB.SendMessage("comp1", m)
	if err != nil {
	}

	fmt.Println(response)

	compA.Broadcast(m)

}
