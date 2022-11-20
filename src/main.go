package main

import (
	"bitTorrent-client/src/torrentFile"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

func main() {
	// inPath := os.Args[1]
	inPath := "debian-11.5.0-amd64-netinst.iso.torrent"
	tf, err := torrentFile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}

	getClients(tf)

	// start by writing a funtion that take torrent file
	// as input and output list of clients

	// err = tf.DownloadToFile("dummy")
	// if err != nil {
	// 	log.Fatal(err)
	// }

}

// "announce" section, specifies the URL of the tracker
// we have to tell tracker we're joining as client
// get base url from torrentFile Announce
// add query params

const Port uint = 6881

type bencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func getClients(tf torrentFile.TorrentFile) {
	fmt.Println(tf.Announce)
	baseUrl, err := url.Parse(tf.Announce)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	baseUrl.RawQuery = url.Values{
		"info_hash":  []string{string(tf.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(tf.Length)},
	}.Encode()

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(baseUrl.String())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Body)
	// result, err := json.Marshal(resp.Body)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	trackerResp := bencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, &trackerResp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	result, _ := json.Marshal(trackerResp)

	fmt.Println(string(result))

}
