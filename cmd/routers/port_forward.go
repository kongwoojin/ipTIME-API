package routers

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func GetPortForwardList(client *http.Client, router *structs.Router) []structs.PortForward {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	req, err := http.NewRequest("GET", baseURL+routerPortForwardList, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", baseURL+routerLogin)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if strings.Contains(html.Text(), "/sess-bin/login_session.cgi") || strings.Contains(html.Text(), "top.location = \"/\";") {
		fmt.Println("Not logged in!")
		return nil
	}

	var portForwardList []structs.PortForward

	table := html.Find("table")
	tbody := table.Find("tbody")

	tbody.Find("tr.pf_tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		id, _ := s.Attr("id")
		if strings.Compare(id, "_-guard_line") == 0 {
			return false
		}

		var portForward = &structs.PortForward{}
		portForward.Name = s.Find("td:nth-child(2) > a:nth-child(1)").Text()
		portForward.IP = s.Find("td:nth-child(3) > span:nth-child(1)").Text()

		portRe := regexp.MustCompile("[0-9]+")
		protocolRe := regexp.MustCompile("(TCPUDP|TCP|UDP|GRE)")

		externalPort := s.Find("td:nth-child(4) > span:nth-child(1)").Text()

		portForward.Protocol = strings.ToLower(protocolRe.FindString(externalPort))

		externalPorts := portRe.FindAllString(externalPort, -1)

		if len(externalPorts) == 2 {
			portForward.ExternalPortStart, _ = strconv.Atoi(externalPorts[0])
			portForward.ExternalPortEnd, _ = strconv.Atoi(externalPorts[1])
		} else {
			portForward.ExternalPortStart, _ = strconv.Atoi(externalPorts[0])
			portForward.ExternalPortEnd = portForward.ExternalPortStart
		}

		internalPort := s.Find("td:nth-child(5) > span:nth-child(1)").Text()
		internalPorts := portRe.FindAllString(internalPort, -1)

		if len(internalPorts) == 2 {
			portForward.InternalPortStart, _ = strconv.Atoi(internalPorts[0])
			portForward.InternalPortEnd, _ = strconv.Atoi(internalPorts[1])
		} else {
			portForward.InternalPortStart, _ = strconv.Atoi(internalPorts[0])
			portForward.InternalPortEnd = portForward.InternalPortStart
		}

		portForwardList = append(portForwardList, *portForward)
		return true
	})

	return portForwardList
}

func AddPortForward(client *http.Client, router *structs.Router, portForward *structs.PortForward) bool {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	switch strings.ToLower(portForward.Protocol) {
	case "tcp":
	case "udp":
	case "tcpudp":
	case "gre":
		break
	default:
		return false
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
		log.Fatal(err)
	}
	req.Header.Set("Referer", baseURL+routerLoginSession)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if strings.Contains(string(data), "onClickedPFRule") {
		return false
	}

	return true
}

func RemovePortForward(client *http.Client, router *structs.Router, portForward *structs.PortForward) {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	params := url.Values{
		"tmenu": []string{"iframe"}, "smenu": []string{"user_portforward"}, "act": []string{"del"},
		"view_mode": []string{"user"}, "mode": []string{"user"},
		"trigger_protocol": []string{""}, "trigger_sport": []string{""}, "trigger_eport": []string{""},
		"forward_ports": []string{""}, "forward_protocol": []string{""},
		"disabled": []string{""}, "priority": []string{""}, "old_priority": []string{""},
		"int_sport": []string{""}, "int_eport": []string{""}, "ext_sport": []string{""}, "ext_eport": []string{""},
		"internal_ip": []string{""}, "protocol": []string{""},

		"delcheck": []string{portForward.Name},
	}

	req, err := http.NewRequest("POST", baseURL+routerRoot, bytes.NewBufferString(params.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Referer", baseURL+routerPortForwardList)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
}
