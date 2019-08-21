package notify

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func SendDingTalk(title, message, robotUrl string, proxyUrl *url.URL) string {
	requestBody := fmt.Sprintf(`{"msgtype": "text","text": {"content": "%s\n\n%s"}}`, title, message)
	jsonStr := []byte(requestBody)

	req, _ := http.NewRequest("POST", robotUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	if proxyUrl != nil {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy:           http.ProxyURL(proxyUrl),
			},
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
	return string(body)
}
