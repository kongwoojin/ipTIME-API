package routers

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetPortForwardList(client *http.Client, router *structs.Router) []structs.PortForward {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	req, err := http.NewRequest("GET", baseURL+routerPortForwardRulesDownload, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", baseURL+routerPortForwardRestore)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(string(data), "/sess-bin/login_session.cgi") || strings.Contains(string(data), "top.location = \"/\";") {
		fmt.Println("Not logged in!")
		return nil
	}

	var portForwardList []structs.PortForward

	portForward := &structs.PortForward{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[") {
			portForward.Name = scanner.Text()[1 : len(scanner.Text())-1]
		} else if strings.Contains(scanner.Text(), "=") {
			s := strings.Split(scanner.Text(), "=")
			key, value := strings.Trim(s[0], " "), strings.Trim(s[1], " ")
			switch key {
			case "enable":
				portForward.Enabled = value == "1"
			case "protocol":
				portForward.Protocol = value
			case "extern_port":
				if strings.Contains(value, "-") {
					s := strings.Split(value, "-")
					portForward.ExternalPortStart, _ = strconv.Atoi(s[0])
					portForward.ExternalPortEnd, _ = strconv.Atoi(s[1])
				} else {
					portForward.ExternalPortStart, _ = strconv.Atoi(value)
					portForward.ExternalPortEnd = portForward.ExternalPortStart
				}
			case "local_port":
				if strings.Contains(value, "-") {
					s := strings.Split(value, "-")
					portForward.InternalPortStart, _ = strconv.Atoi(s[0])
					portForward.InternalPortEnd, _ = strconv.Atoi(s[1])
				} else {
					portForward.InternalPortStart, _ = strconv.Atoi(value)
					portForward.InternalPortEnd = portForward.InternalPortStart
				}
			case "local_ip":
				portForward.IP = value
				portForwardList = append(portForwardList, *portForward)
				portForward = &structs.PortForward{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	return portForwardList
}

func AddPortForward(client *http.Client, router *structs.Router, portForward *structs.PortForward) (bool, error) {
	if checkPortForwardExist(client, router, portForward.Name) {
		return false, fmt.Errorf("portforward rule \"%s\" already exist", portForward.Name)
	}

	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	switch strings.ToLower(portForward.Protocol) {
	case "tcp":
	case "udp":
	case "tcpudp":
	case "gre":
		break
	default:
		return false, fmt.Errorf("Unknown protocol: %s", portForward.Protocol)
	}

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"user_portforward"}, "act": []string{"add"},
		"view_mode": []string{"user"}, "mode": []string{"user"},
		"trigger_protocol": []string{""}, "trigger_sport": []string{""}, "trigger_eport": []string{""},
		"forward_ports": []string{""}, "forward_protocol": []string{""},
		"disabled": []string{"0"}, "priority": []string{""}, "old_priority": []string{""},

		"name":        []string{portForward.Name},
		"int_sport":   []string{strconv.Itoa(portForward.InternalPortStart)},
		"int_eport":   []string{strconv.Itoa(portForward.InternalPortEnd)},
		"ext_sport":   []string{strconv.Itoa(portForward.ExternalPortStart)},
		"ext_eport":   []string{strconv.Itoa(portForward.ExternalPortEnd)},
		"internal_ip": []string{portForward.IP},
		"protocol":    []string{strings.ToLower(portForward.Protocol)},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Referer", baseURL+routerLoginSession)
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

	if strings.Contains(string(data), portForward.Name) {
		return true, nil
	} else {
		return false, fmt.Errorf("failed to add portforward rule \"%s\"", portForward.Name)
	}
}

func RemovePortForward(client *http.Client, router *structs.Router, portForward *structs.PortForward) (bool, error) {
	if !checkPortForwardExist(client, router, portForward.Name) {
		return false, fmt.Errorf("portforward rule \"%s\" cannot be found", portForward.Name)
	}

	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"user_portforward"}, "act": []string{"del"},
		"view_mode": []string{"user"}, "mode": []string{""},
		"trigger_protocol": []string{""}, "trigger_sport": []string{""}, "trigger_eport": []string{""},
		"forward_ports": []string{""}, "forward_protocol": []string{""},
		"disabled": []string{""}, "priority": []string{""}, "old_priority": []string{""},
		"int_sport": []string{""}, "int_eport": []string{""}, "ext_sport": []string{""}, "ext_eport": []string{""},
		"internal_ip": []string{""}, "protocol": []string{""},

		"delcheck": []string{portForward.Name},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Referer", baseURL+routerPortForwardList)
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

	if strings.Contains(string(data), portForward.Name) {
		return false, fmt.Errorf("failed to remove portforward rule \"%s\"", portForward.Name)
	} else {
		return true, nil
	}
}

func checkPortForwardExist(client *http.Client, router *structs.Router, portForwardName string) bool {
	for _, portForwardItem := range GetPortForwardList(client, router) {
		if strings.Compare(portForwardItem.Name, portForwardName) == 0 {
			return true
		}
	}
	return false
}
