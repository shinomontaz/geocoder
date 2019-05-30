package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/shinomontaz/geocoder/service/proxy"
)

type ILogger interface {
	Log(message string)
}

type IProxy interface {
	Take(limit int) *http.Client
	UpdateWindows()
}

type IGeocoder interface {
	Query(address string, client *http.Client) (x, y float64, err error)
}

type Coder struct {
	Coder IGeocoder
	Limit int
}

type RpcServer struct {
	px         *proxy.Proxy
	logger     ILogger
	proxy      IProxy
	coders     map[string]*Coder
	coderUsage map[string]int
	mu         sync.RWMutex
}

func (rs *RpcServer) Add(name string, geocoder IGeocoder, limit int) {
	rs.coders[name] = &Coder{Coder: geocoder, Limit: limit}
}

func (rs *RpcServer) Del(coder string) {
	delete(rs.coders, coder)
}

func (rs *RpcServer) TakeCoder() *Coder {
	randCoder := rand.Intn(len(rs.coders))
	keys := make([]string, 0, len(rs.coders))
	for key, _ := range rs.coders {
		keys = append(keys, key)
	}
	return rs.coders[keys[randCoder]]
}

type Res struct {
	X       float64 `json:"lat"`
	Y       float64 `json:"long"`
	Success string  `json:"success"`
}

type GeocodeQuery struct {
	Address string `json:"address"`
}

type GeocodeResponse struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
	err  error
}

func (rs *RpcServer) Geocode(r *http.Request, req *GeocodeQuery, resp *GeocodeResponse) error { // это запускается в го-процедуре
	// select geocoder
	// take proxy
	// geocode and parse result

	coder := rs.TakeCoder()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	ready := make(chan GeocodeResponse)

	go func() { // Это нужно чтоб выйти по context.Timeout
		proxyClient := rs.proxy.Take(coder.Limit)
		Lat, Long, err := coder.Coder.Query(req.Address, proxyClient)
		ready <- GeocodeResponse{Lat: Lat, Long: Long, err: err}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-ready:
		if result.err != nil {
			return result.err
		}
		*resp = result
	}

	return nil
}

func (rs *RpcServer) Start() {
	//	longTick := time.NewTicker(31 * time.Minute)
	longTick := time.NewTicker(1 * time.Minute)

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-longTick.C:
				rs.longTick()
			case <-quit:
				longTick.Stop()
				return
			}
		}
	}()
}

func (rs *RpcServer) longTick() {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	fmt.Println("longTick")
	rs.proxy.UpdateWindows()
}
