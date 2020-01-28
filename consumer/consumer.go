package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8090", "http service address")

// msg is our in memory store of data we get from poducer and send to client
var msg struct {
	Data string `json:"data,omitempty"`
}

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// streamHandler is constantly printing data from msg to SSE connection to client
func streamHandler(w http.ResponseWriter, r *http.Request) {

	var buf = new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(msg)

	w.Header().Set("Content-Type", "text/event-stream") // SSE specific header
	w.Header().Set("Cache-Control", "no-cache")         // make sure data is always fresh
	w.Header().Set("Connection", "keep-alive")          // some say not necessary, but keeps browser from closing connection
	fmt.Fprintf(w, "data: %v\n\n", buf)                 // writing our json encoded message to stream
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// reader is called by wsHandler to read data coming over ws connection from
// producer into msg, our in memory data store to write to SSE stream to user
func reader(conn *websocket.Conn) {
	err := conn.ReadJSON(&msg) // read our json from ws connection into msg
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("read message:", msg)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // ordinarily we would check origin against approved list
	ws, err := upgrader.Upgrade(w, r, nil)                            // our ws connection is initialized

	if err != nil {
		log.Println("wsHandler", err)
	}
	log.Println("Client Connected")

	reader(ws) // we pass our ws connection to reader
}
