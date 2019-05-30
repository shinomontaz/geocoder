package proxy

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type pServer struct {
	Url    url.URL
	Auth   string
	client *http.Client
	mu     sync.Mutex
	Usage  []int64
}

func NewPServer(proxyUrl *url.URL, proxyAuth string) *pServer {
	newPServer := &pServer{
		Url:  *proxyUrl,
		Auth: proxyAuth,
		client: &http.Client{
			Timeout:   time.Duration(10 * time.Second),
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
		},
	}

	header := http.Header{}
	auth := proxyAuth
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	header.Add("Proxy-Authorization", basicAuth)

	newPServer.client.Transport = &http.Transport{
		Proxy:              http.ProxyURL(proxyUrl),
		ProxyConnectHeader: header,
	}

	return newPServer
}

func (p *pServer) updateWindow() {
	now := time.Now().Unix()

	fmt.Println("update window", p.Url.Host, len(p.Usage))
	if len(p.Usage) > 10 {
		idx := int(float64(len(p.Usage)) / 33.0)
		if p.Usage[idx] < (now - 24*60*60) {
			p.Usage = p.Usage[idx:]
		}
	}
}

func (p *pServer) UpdateWindow() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.updateWindow()
}

func (p *pServer) GetClient() *http.Client {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.updateWindow()

	now := time.Now().Unix()
	p.Usage = append(p.Usage, now)

	return p.client
}
