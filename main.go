package main

import (
	"log"
	"net/http"

	"geocoder/service/geocoder"
	"geocoder/service/logger"
	"geocoder/service/proxy"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
)

/*
Пример запроса на адрес 127.0.0.1:1111:
{
"jsonrpc":"2.0",
"method":"RpcServer.Geocode",
"params":[{
	"address": "г.Москва, ул. Грина, 42"
}],
"id": 1
}

а в ответ -
*/

func main() {
	logService := logger.New(slackClient, make(chan string, 100), cfg.SlackChannel)
	logService.Start()

	proxyService := proxy.New()

	for _, entry := range cfg.List {
		proxyService.Add(entry.Url, entry.Auth)
	}

	coderYandex := geocoder.Yandex()
	//	coderSputnik := geocoder.Sputnik()

	handler := &RpcServer{logger: logService, proxy: proxyService, coders: make(map[string]*Coder, 0)}

	handler.Start()

	handler.Add("yandex", coderYandex, 20000)
	//	handler.Add("sputnik", coderSputnik, 20000)

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	err := s.RegisterService(handler, "")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", s)

	log.Fatal(http.ListenAndServe(":1111", nil)) // здесь каждое соединение обрабатывается в отдельной го-процедуре
}
