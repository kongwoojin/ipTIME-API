package routers

import (
	"fmt"
	"github.com/kongwoojin/ipTIME-API/cmd/structs"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func Login(client *http.Client, router *structs.Router) bool {
	var baseURL = "http://" + router.Host + ":" + fmt.Sprint(router.Port) + "/sess-bin/"

	params := url.Values{
		"init_status": []string{"1"}, "captcha_on": []string{"0"}, "captcha_file": []string{},
		"default_passwd": []string{}, "username": {router.Username}, "passwd": {router.Password}, "captcha_code": []string{},
	}

	postData := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", baseURL+routerLoginHandler, postData)
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

	if strings.Contains(string(data), "/sess-bin/login_session.cgi") || strings.Contains(string(data), "top.location = \"/\";") {
		return false
	} else {
		regex := regexp.MustCompile(`setCookie\('([^']+)'\);`)

		match := regex.FindStringSubmatch(string(data))
		if match != nil {
			extractedValue := match[1]

			var cookies []*http.Cookie

			cookie := &http.Cookie{
				Name:  "efm_session_id",
				Value: extractedValue,
				Path:  "/",
			}

			cookies = append(cookies, cookie)

			parsedUrl, err := url.Parse(baseURL)

			if err != nil {
				log.Fatal(err)
			}

			client.Jar.SetCookies(parsedUrl, cookies)
		}
	}
	return true
}
