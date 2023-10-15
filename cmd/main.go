package main

import (
	"github.com/kongwoojin/ipTIME-API/cmd/routers"
	"log"
)

func main() {
	router := RouterInfo()
	client := CreateClient()

	if routers.Login(client, router) {
		routers.RouterStatus(client, router)
		routers.GetConnectedClientList(client, router)
	} else {
		log.Fatal("Login failed")
	}
}
