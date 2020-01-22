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

var msg struct {
	Data string `json:"data,omitempty"`
}

// var msg = new(Msg)

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("client")))
	// http.Handle("/echo", http.HandlerFunc(echoHandler))
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {

	var buf = new(bytes.Buffer)

	enc := json.NewEncoder(buf)
	enc.Encode(msg)
	fmt.Printf("data: %v\n", buf.String())

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	fmt.Fprintf(w, "data: %v\n\n", buf)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	// read in a message
	err := conn.ReadJSON(&msg)
	if err != nil {
		log.Println("lll", err)
		return
	}

	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("read message:", msg)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("wsHandler", err)
	}
	log.Println("Client Connected")

	reader(ws)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	return
}
