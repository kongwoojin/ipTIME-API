package routers

import (
	"bytes"
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func AddWOL(client *http.Client, router *structs.Router, macAddress string, name string) (bool, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	s := strings.Split(macAddress, ":")

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"expertconfwollist"}, "act": []string{"add"},
		"nomore": []string{"0"},

		"hw1": []string{s[0]}, "hw2": []string{s[1]}, "hw3": []string{s[2]}, "hw4": []string{s[3]}, "hw5": []string{s[4]}, "hw6": []string{s[5]},
		"pcname": []string{name},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Referer", baseURL+routerRoot)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if strings.Contains(string(data), strings.ReplaceAll(macAddress, ":", "-")) {
		return true, nil
	} else {
		return false, fmt.Errorf("failed to add %s to Wake on Lan list", macAddress)
	}
}

func RemoveWOL(client *http.Client, router *structs.Router, macAddress string) (bool, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"expertconfwollist"}, "act": []string{"del"},
		"nomore": []string{"0"},

		"delchk": []string{macAddress},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Referer", baseURL+routerWOLList)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if strings.Contains(string(data), strings.ReplaceAll(macAddress, ":", "-")) {
		return false, fmt.Errorf("failed to remove %s from Wake on Lan list", macAddress)
	} else {
		return true, nil
	}
}

func Wake(client *http.Client, router *structs.Router, macAddress string) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"expertconfwollist"}, "act": []string{"wake"},
		"nomore": []string{"0"},

		"delchk": []string{macAddress},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", baseURL+routerWOLList)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
