package main

import (
	"bytes"
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
	http.Handle("/", http.HandlerFunc(indexHandler))
	http.HandleFunc("/ws", wsEndpoint)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("this will be sent to consumer:", string(body[:]))

	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	resp, err := http.Post("http://127.0.0.1:8090/ws", "text/html", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	fmt.Println("this was returned from consumer:", string(body[:]))

	return
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
	// send message
	err = c.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		fmt.Println("writeMesssage", err)
	}
	// receive message
	_, message, err := c.ReadMessage()
	if err != nil {
		fmt.Println("readMesssage", err)
	}
	fmt.Println("returnMesssage", message)
}
