package main

import (
	"encoding/json"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"os"
)

func RouterInfo() *structs.Router {
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
