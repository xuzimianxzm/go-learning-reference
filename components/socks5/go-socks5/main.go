package main

import (
	"fmt"
	"github.com/armon/go-socks5"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// Create a SOCKS5 server
	conf := &socks5.Config{}
	go startSocks5Server(conf)

	time.Sleep(1500)
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	// set our socks5 as the dialer
	httpTransport.DialContext = conf.Dial

	if resp, err := httpClient.Get("https://wtfismyip.com/json"); err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}
}

func startSocks5Server(conf *socks5.Config) {
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "127.0.0.1:8000"); err != nil {
		log.Fatalln(err)
	}
}
