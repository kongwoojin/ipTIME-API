package main

import (
	"encoding/json"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
)

func main() {
	router := routerInfo()
	client := createClient()

	if login(client, router) {
		routerStatus(client, router)
	} else {
		log.Fatal("Login failed")
	}
}

func routerInfo() *structs.Router {
	jsonFile, err := os.Open("info.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var router *structs.Router

	json.Unmarshal(byteValue, &router)

	return router
}

func createClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}

	client := &http.Client{
		Jar: jar,
	}

	return client
}
