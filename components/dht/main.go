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
	getFileInfoForURI(torrent, client)

	go startToDownload(torrent, client)

	readDownloadFileSpeed(torrent)

	log.Print("torrent downloaded")
}

func readDownloadFileSpeed(torrent *torrent.Torrent) {
	timeTickerChan := time.Tick(time.Second * 2)
	i := 0
	for {
		if i > 100 {
			break
		}
		i++
		torrentStats := torrent.Stats()
		log.Print("downloaded BytesReadData ", torrentStats.BytesReadData.String())
		log.Print("downloaded BytesWrittenData ", torrentStats.BytesWrittenData.String())
		<-timeTickerChan
	}
}

func startToDownload(torrent *torrent.Torrent, client *torrent.Client) {
	torrent.DownloadAll()
	client.WaitAll()
}

func getFileInfoForURI(torrent *torrent.Torrent, client *torrent.Client) {
	infoDone := <-torrent.GotInfo()
	for _, file := range torrent.Files() {
		log.Print("torrent Files ", file.Path())
		log.Print("torrent Length ", file.Length())
		log.Print("torrent BytesCompleted ", file.BytesCompleted())
		log.Print("torrent FileInfo ", file.FileInfo())
	}

	clientStates := client.ConnStats()
	log.Print("client ConnStats BytesWritten ", clientStates.BytesWritten.String())
	log.Print("client ConnStats BytesReadData ", clientStates.BytesReadData.String())
	log.Print("have completed download bytes ", torrent.BytesCompleted())

	log.Println("got info Done:", infoDone)
	info := torrent.Info()
	log.Println("Info Name", info.Name)
	log.Println("Info Length", info.Length)
	log.Println("Info Pieces", len(info.Pieces))

	metainfo := torrent.Metainfo()
	log.Println("Metainfo InfoBytes ", len(metainfo.InfoBytes))
	log.Println("Metainfo CreatedBy ", metainfo.CreatedBy)
	log.Println("Metainfo CreatedBy ", metainfo.Comment)
}
