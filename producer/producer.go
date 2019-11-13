package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var addrCon = flag.String("addrCon", ":8090", "http service address")

func main() {
	flag.Parse()
	http.HandleFunc("/", wsEndpoint)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("this will be sent to consumer:", string(body[:]))

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	u := url.URL{Scheme: "ws", Host: *addrCon, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dialerError", err)
		// handle error
	}

	type Msg struct {
		Data string `json:"data,omitempty"`
	}

	msg := Msg{Data: string(body)}

	var buf = new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.Encode(msg)

	if err != nil {
		fmt.Println("error:", err)
	}

	var expect struct {
		Data string
	}
	expect.Data = string(body)

	// send message
	err = c.WriteJSON(&expect)
	if err != nil {
		fmt.Println("writeJSON", err)
	}
	// receive message
	// _, message, err := c.ReadMessage()
	// if err != nil {
	// 	fmt.Println("readMesssage", err)
	// }
	return
}
