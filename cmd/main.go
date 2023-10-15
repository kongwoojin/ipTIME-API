package main

import (
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/routers"
	"log"
)

// main function for testing API call
func main() {
	router := RouterInfo()
	client := CreateClient()

	if routers.Login(client, router) {
		// Check Router status
		fmt.Printf("%+v\n", *routers.RouterStatus(client, router))

		// Check connected clients
		for i, v := range routers.GetConnectedClientList(client, router) {
			fmt.Printf("%d: %+v\n", i, v)
		}
	} else {
		log.Fatal("Login failed")
	}
}
