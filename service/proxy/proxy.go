package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type Stats struct {
	requests int
	success  int
	timeout  int
	Usage    map[string][]int64
}

type Proxy struct {
	mu          sync.Mutex
	currProxies []*pServer // store response from proxy list site
	stats       *Stats
}

func New() *Proxy {
	return &Proxy{
		stats:       &Stats{},
		currProxies: make([]*pServer, 0),
	}
}

func (p *Proxy) Take(limit int) *http.Client {
	p.mu.Lock()
	defer p.mu.Unlock()

	sort.Slice(p.currProxies, func(i, j int) bool {
		return len(p.currProxies[i].Usage) < len(p.currProxies[j].Usage)
	})

	if len(p.currProxies[0].Usage) > limit {
		return &http.Client{Timeout: time.Duration(10 * time.Second)}
	}
	for _, us := range p.currProxies[0].Usage {
		fmt.Println(time.Unix(us, 0).Format(time.RFC3339))
	}
	fmt.Println("take proxy from list", p.currProxies[0].Url.Host)
	return p.currProxies[0].GetClient()
}

func (p *Proxy) Add(proxyAddress, proxyAuth string) {

	if strings.Index(proxyAddress, "http://") == -1 {
		proxyAddress = fmt.Sprintf("http://%s", proxyAddress)
	}

	proxyUrl, err := url.Parse(proxyAddress)
	if err != nil {
		fmt.Println("parsing proxy url:", err)
	}

	newPServer := NewPServer(proxyUrl, proxyAuth)

	for _, existedPserver := range p.currProxies {
		if existedPserver.Url == newPServer.Url {
			return
		}
	}

	p.currProxies = append(p.currProxies, newPServer)
}

func (p *Proxy) UpdateWindows() {
	for _, pr := range p.currProxies {
		go pr.UpdateWindow() // no need to wait and do it one by one
	}
}
