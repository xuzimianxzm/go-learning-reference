package main

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	_ "io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

// used socks5 dialer
// Notes: A socks5 server is required when the program is running before,
// for example: go-socks/main.go startSocks5Server()
func main() {
	// create a socks5 dialer
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:8000", nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}

	// set our socks5 as the dialer
	dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialer.Dial(network, address)
	}
	httpTransport := &http.Transport{
		DialContext:       dialContext,
		DisableKeepAlives: true,
	}
	httpClient := &http.Client{Transport: httpTransport}
	if resp, err := httpClient.Get("https://wtfismyip.com/json"); err != nil {
		panic(err)
	} else {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%s\n", body)
	}
}
