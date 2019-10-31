package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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

func SendDingTalkNew(title, message, robotUrl, secret string, proxyUrl *url.URL) string {
	timestamp := fmt.Sprintf("%d000", time.Now().Unix())
	sign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sign))

	signB64 := base64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))

	v := url.Values{}
	v.Add("sign", signB64)
	signUrlEncode := v.Encode()
	postUrl := fmt.Sprintf("%s&timestamp=%s&%s", robotUrl, timestamp, signUrlEncode)

	requestBody := fmt.Sprintf(`{"msgtype": "text","text": {"content": "%s\n\n%s"}}`, title, message)
	jsonStr := []byte(requestBody)

	req, _ := http.NewRequest("POST", postUrl, bytes.NewBuffer(jsonStr))
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
