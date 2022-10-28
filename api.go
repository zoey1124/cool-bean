package main

import (
	"fmt"
	"os"

	"github.com/cs161-staff/project2-starter-code/client"
)

func main() {
	fmt.Println(os.Args[1])
	if len(os.Args) < 2 {
		fmt.Println("Missing parameter")
		return
	}
	command := os.Args[1]
	switch command {
	case "InitUser":
		user, _ := client.InitUser("Alice", "password")
		fmt.Println(user.Username)
	case "GetUser":
		// ...
	}
}
