package main

import (
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/routers"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
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

		// Check port forward list
		for i, v := range routers.GetPortForwardList(client, router) {
			fmt.Printf("%d: %+v\n", i, v)
		}

		// Check add port forward
		routers.AddPortForward(client, router, &structs.PortForward{
			Name:              "portforward_test",
			IP:                "192.168.0.253",
			Protocol:          "TCP",
			InternalPortStart: 80,
			InternalPortEnd:   80,
			ExternalPortStart: 80,
			ExternalPortEnd:   80,
		})

		routers.RemovePortForward(client, router, &structs.PortForward{
			Name: "portforward_test",
		})
	} else {
		log.Fatal("Login failed")
	}
}
