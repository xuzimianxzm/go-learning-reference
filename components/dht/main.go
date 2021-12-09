package main

import (
	"github.com/anacrolix/torrent"
	"log"
	"time"
)

func main() {

	config := torrent.NewDefaultClientConfig()
	config.DataDir = "/Users/zimian.xu/Downloads/"
	config.ListenPort = 42070 // Listen default UDP Port:42069

	client, err := torrent.NewClient(config)
	if err != nil {
		log.Print("--- NewClient Error :", err)
		return
	}
	defer client.Close()
	torrent, _ := client.AddMagnet("magnet:?xt=urn:btih:ZOCMZQIPFFW7OLLMIC5HUB6BPCSDEOQU")
	info := <-torrent.GotInfo()
	log.Print("--- GotInfo ---:", info)
	log.Println("--- Info ---")
	log.Println("Info Name", torrent.Info().Name)
	log.Println("Info Length", torrent.Info().Length)
	torrent.DownloadAll()
	client.WaitAll()
	time.Sleep(1000 * 10)
	log.Print("torrent downloaded")
}
