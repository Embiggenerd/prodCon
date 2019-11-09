package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8090", "http service address")
var msg []byte

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.Handle("/echo", http.HandlerFunc(echoHandler))
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/ws", wsHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	fmt.Fprintf(w, "data: %v\n\n", string(msg))
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	fmt.Println("bowdy", string(body))
	msg = body
	// w.Write(body)
	return
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		fmt.Println("received msg", string(p))
		msg = p
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}
func wsHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("zzz", err)
	}
	log.Println("Client Connected")

	reader(ws)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// http.Handle("/", http.FileServer(http.Dir("client")))
	fmt.Println("msg", string(msg))
	w.Write(msg)
	return
}
