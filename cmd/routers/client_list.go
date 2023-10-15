package routers

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"log"
	"net/http"
	"strings"
)

func GetConnectedClientList(client *http.Client, router *structs.Router) []structs.Client {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	req, err := http.NewRequest("GET", baseURL+routerLanPCInfoStatus, nil)
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

	var connectedClients []structs.Client

	table := html.Find("table.lansetup_main_table")
	tbody := table.Find("tbody")

	tbody.Find("tr.lansetup_main_tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
		var client = &structs.Client{}

		ip := strings.Trim(s.Find("td:nth-child(1) > span:nth-child(1)").Text(), " ")

		if ip == "" {
			return false
		}

		client.IP = ip
		client.Mac = strings.Trim(s.Find("td:nth-child(2) > span:nth-child(1)").Text(), " ")
		client.Hostname = strings.Trim(s.Find("td:nth-child(3) > span:nth-child(1)").Text(), " ")

		connectionType := strings.Split(s.Find("td:nth-child(4) > span:nth-child(1)").Text(), ":")

		switch connectionType[0] {
		case "유선":
			client.ConnectionType = structs.Wired
		case "무선":
			client.ConnectionType = structs.Wireless
		default:
			client.ConnectionType = structs.UnknownConnectionType
		}

		if strings.Compare(connectionType[1], "수동할당") == 0 {
			client.IsManuallyAssigned = true
		} else {
			client.IsManuallyAssigned = false
		}

		connectedClients = append(connectedClients, *client)
		return true
	})

	return connectedClients
}
