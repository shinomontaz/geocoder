package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	var jsonStr = []byte(`{
		"jsonrpc":"2.0",
		"method":"RpcServer.Geocode",
		"params":[{
			"Address": "г.Москва, ул. Грина, 42"
		}],
		"id": 1
		}`)
	//192.168.1.1:1111
	req, err := http.NewRequest("POST", "http://192.168.1.4:1111", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp1, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp1.Body)
	fmt.Println(string(body))

	resp2, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ = ioutil.ReadAll(resp2.Body)
	fmt.Println(string(body))

	// proxyUrl, _ := url.Parse("http://46.232.62.89:7598")

	// header := http.Header{}

	// auth := "user19494:pn5fdu"
	// basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	// header.Add("Proxy-Authorization", basicAuth)

	// client := &http.Client{Transport: &http.Transport{
	// 	Proxy:              http.ProxyURL(proxyUrl),
	// 	ProxyConnectHeader: header,
	// }}

	// resp, err := client.Get("https://geocode-maps.yandex.ru/1.x/?format=json&geocode=Домодедовская, дом 37, к.2")

	// if err != nil {
	// 	panic(err)
	// }

	// result, _ := ioutil.ReadAll(resp.Body)
	// resp.Body.Close()

	// fmt.Println(string(result))

	//	делать одновременные запросы и очень быстро - по 10 штук за раз
}
