package routers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func RouterStatus(client *http.Client, router *structs.Router) *structs.RouterStatus {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	req, err := http.NewRequest("GET", baseURL+routerSystemInfoStatus, nil)
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

	routerStatus := &structs.RouterStatus{}

	table := html.Find("table.menu_content_list_noline_table")
	tbody := table.Find("tbody")

	// Get Router's Internet Connection Status
	parsedInternetStatus := strings.Trim(tbody.Find("tr.odd:nth-child(2) > td:nth-child(2)").Text(), " ")

	if strings.Compare(parsedInternetStatus, "인터넷에 정상적으로 연결됨") == 0 {
		routerStatus.IsInternetConnected = true
	} else {
		routerStatus.IsInternetConnected = false
	}

	// Get Router's IP Type
	parsedIpType := strings.Trim(tbody.Find("tr.even:nth-child(3) > td:nth-child(2) > span:nth-child(1)").Text(), " ")

	switch parsedIpType {
	case "동적 IP 연결":
		routerStatus.IpType = structs.DynamicIp
	case "고정 IP 연결":
		routerStatus.IpType = structs.StaticIp
	default:
		routerStatus.IpType = structs.Unknown
	}

	// Get Router's IP
	ip := strings.Trim(tbody.Find("td.td_item:nth-child(4) > span:nth-child(1)").Text(), " ")
	routerStatus.Ip = ip

	// Get Router's 5Ghz Wifi Info
	fiveGhzInfo := structs.WifiInfo{}

	parsed5ghzStatus := strings.Trim(tbody.Find("tr.odd:nth-child(10) > td:nth-child(2)").Text(), " ")
	if strings.Compare(parsed5ghzStatus, "중단됨") == 0 {
		fiveGhzInfo.IsOn = false
	} else {
		fiveGhzInfo.IsOn = true
	}

	fiveGhzInfo.Ssid = strings.Trim(tbody.Find("tr.even:nth-child(11) > td:nth-child(2)").Text(), " ")

	routerStatus.Fiveghz = fiveGhzInfo

	// Get Router's 2.4Ghz Wifi Info
	twoGhzInfo := structs.WifiInfo{}

	parsed2ghzStatus := strings.Trim(tbody.Find("tr.odd:nth-child(16) > td:nth-child(2)").Text(), " ")
	if strings.Compare(parsed2ghzStatus, "중단됨") == 0 {
		twoGhzInfo.IsOn = false
	} else {
		twoGhzInfo.IsOn = true
	}

	twoGhzInfo.Ssid = strings.Trim(tbody.Find("tr.even:nth-child(17) > td:nth-child(2)").Text(), " ")
	routerStatus.Twoghz = twoGhzInfo

	// Get Router's Firmware Version
	routerStatus.FirmwareVersion = strings.Trim(tbody.Find("tr.odd:nth-child(22) > td:nth-child(2)").Text(), " ")

	// Get Router's Remote Access Status
	parsedRemoteAccessStatus := strings.Trim(tbody.Find("tr.even:nth-child(23) > td:nth-child(2)").Text(), " ")
	if strings.Compare(parsedRemoteAccessStatus, "원격 관리 포트가 설정되어 있지 않음") == 0 {
		routerStatus.RemoteAccess = false
	} else {
		routerStatus.RemoteAccess = true
		re := regexp.MustCompile("[0-9]+")

		port, err := strconv.Atoi(re.FindString(parsedRemoteAccessStatus))
		if err != nil {
			routerStatus.RemoteAccess = false
		} else {
			routerStatus.RemoteAccessPort = port
		}
	}

	// Get Router's Uptime
	parsedUptime := strings.Trim(tbody.Find("tr.even:nth-child(25) > td:nth-child(2)").Text(), " ")
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(parsedUptime, -1)

	uptime := structs.Uptime{}

	uptime.Days, err = strconv.Atoi(matches[0])
	if err != nil {
		log.Fatal(err)
	}

	uptime.Hours, err = strconv.Atoi(matches[1])
	if err != nil {
		log.Fatal(err)
	}

	uptime.Minutes, err = strconv.Atoi(matches[2])
	if err != nil {
		log.Fatal(err)
	}

	uptime.Seconds, err = strconv.Atoi(matches[3])
	if err != nil {
		log.Fatal(err)
	}

	routerStatus.Uptime = uptime

	return routerStatus
}
