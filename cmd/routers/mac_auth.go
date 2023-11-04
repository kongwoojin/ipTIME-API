package routers

import (
	"bytes"
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/enums"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func ChangeMacAuthMode(client *http.Client, router *structs.Router, policy enums.MacAuthPolicy, frequency enums.WifiFrequency) (bool, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/"

	params := url.Values{
		"tmenu": []string{"wirelessconf"}, "smenu": []string{"macauth"}, "service_name": []string{"macauth"},
		"wlmode": []string{frequency.GetFrequency()}, "idx": []string{"0"}, "mode": []string{"policy"},
		"policy": []string{policy.GetPolicy()},
	}

	req, err := http.NewRequest("POST", baseURL+mobileRouterSubmit, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+mobileRouterMacAuth)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(data), "ok") {
		return true, nil
	} else {
		return false, fmt.Errorf("failed to change mac auth policy")
	}
}

func AddMacAuth(client *http.Client, router *structs.Router, frequency enums.WifiFrequency, macAddress string, desc string) (bool, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/"

	params := url.Values{
		"tmenu": []string{"wirelessconf"}, "smenu": []string{"macauth"}, "service_name": []string{"macauth"},
		"wlmode": []string{frequency.GetFrequency()}, "idx": []string{"0"}, "mode": []string{"register"},
		"mac": []string{macAddress}, "desc": []string{desc},
	}

	req, err := http.NewRequest("POST", baseURL+mobileRouterSubmit, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+mobileRouterMacAuth)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(data), "ok") {
		return true, nil
	} else {
		return false, fmt.Errorf("failed to add mac address")
	}
}

func RemoveMacAuth(client *http.Client, router *structs.Router, frequency enums.WifiFrequency, macAddress string) (bool, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/"

	params := url.Values{
		"tmenu": []string{"wirelessconf"}, "smenu": []string{"macauth"}, "service_name": []string{"macauth"},
		"wlmode": []string{frequency.GetFrequency()}, "idx": []string{"0"}, "mode": []string{"unregister"},
		"delmaccheck": []string{strings.ReplaceAll(macAddress, ":", "-")},
	}

	req, err := http.NewRequest("POST", baseURL+mobileRouterSubmit, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return false, err
	}

	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+mobileRouterMacAuth)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(data), "ok") {
		return true, nil
	} else {
		return false, fmt.Errorf("failed to remove mac address")
	}
}
