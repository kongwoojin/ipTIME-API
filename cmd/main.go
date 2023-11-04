package main

import (
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/enums"
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
		status, err := routers.RouterStatus(client, router)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%+v\n", *status)

		// Check connected clients
		for i, v := range routers.GetConnectedClientList(client, router) {
			fmt.Printf("%d: %+v\n", i, v)
		}

		// Check port forward list
		for i, v := range routers.GetPortForwardList(client, router) {
			fmt.Printf("%d: %+v\n", i, v)
		}

		// Check add port forward
		_, addErr := routers.AddPortForward(client, router, &structs.PortForward{
			Name:              "portforward_test",
			IP:                "192.168.0.253",
			Protocol:          "TCP",
			InternalPortStart: 80,
			InternalPortEnd:   80,
			ExternalPortStart: 80,
			ExternalPortEnd:   80,
		})
		if addErr != nil {
			log.Fatal(addErr)
		}

		_, removeErr := routers.RemovePortForward(client, router, &structs.PortForward{
			Name: "portforward_test",
		})
		if removeErr != nil {
			log.Fatal(removeErr)
		}

		// Check add WOL
		_, addWolErr := routers.AddWOL(client, router, "00:00:00:00:00:00", "test")
		if addWolErr != nil {
			log.Fatal(addWolErr)
		}

		// Check remove WOL
		_, removeWolErr := routers.RemoveWOL(client, router, "00:00:00:00:00:00")
		if removeWolErr != nil {
			log.Fatal(removeWolErr)
		}

		// Send WOL
		routers.Wake(client, router, "00:00:00:00:00:00")

		// Change Mac auth policy
		_, macAuthModeErr := routers.ChangeMacAuthMode(client, router, enums.WhiteList, enums.F5GHZ)
		if macAuthModeErr != nil {
			log.Fatal(macAuthModeErr)
		}

		// Add Mac auth
		_, addMacAuthErr := routers.AddMacAuth(client, router, enums.F5GHZ, "AA:BB:CC:DD:EE:FF", "test")
		if addMacAuthErr != nil {
			log.Fatal(addMacAuthErr)
		}

		// Remove Mac auth
		_, removeMacAuthErr := routers.RemoveMacAuth(client, router, enums.F5GHZ, "AA:BB:CC:DD:EE:FF")
		if removeMacAuthErr != nil {
			log.Fatal(removeMacAuthErr)
		}
	} else {
		log.Fatal("Login failed")
	}
}
