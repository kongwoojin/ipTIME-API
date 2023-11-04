package routers

import (
	"encoding/json"
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"net/http"
)

func RouterStatus(client *http.Client, router *structs.Router) (*structs.RouterStatus, error) {
	routerStatus := &structs.RouterStatus{}

	networkStatus, err := getNetworkStatus(client, router)
	if err != nil {
		return nil, err
	}
	routerStatus.NetworkStatus = *networkStatus

	systemInfo, err := getSystemInfo(client, router)
	if err != nil {
		return nil, err
	}
	routerStatus.SystemInfo = *systemInfo

	return routerStatus, nil
}

func getNetworkStatus(client *http.Client, router *structs.Router) (*structs.NetworkStatus, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/"

	req, err := http.NewRequest("GET", baseURL+mobileRouterStatus, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", baseURL+mobileRouterRoot)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	networkStatus := &structs.NetworkStatus{}

	err = json.Unmarshal(data, &networkStatus)
	if err != nil {
		return nil, err
	}

	return networkStatus, nil
}

func getSystemInfo(client *http.Client, router *structs.Router) (*structs.SystemInfo, error) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/"

	req, err := http.NewRequest("GET", baseURL+mobileRouterSystemInfo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", baseURL+mobileRouterRoot)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	systemInfo := &structs.SystemInfo{}

	err = json.Unmarshal(data, &systemInfo)
	if err != nil {
		return nil, err
	}

	return systemInfo, nil
}
