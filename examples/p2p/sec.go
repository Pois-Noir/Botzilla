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

	_, err := botzilla.RegisterComponent("localhost:6985", "comp2", 6969, amirgay)

	fmt.Println("fffffff")
	if err != nil {
		fmt.Println("There was an error running example, ", err)
		os.Exit(1)
	}

	for true {
	}

}
