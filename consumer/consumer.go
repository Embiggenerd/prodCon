package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8090", "http service address")
var msg []byte

func main() {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("client")))
	http.Handle("/echo", http.HandlerFunc(echoHandler))
	http.HandleFunc("/stream", streamHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// var buf bytes.Buffer
	// enc := json.NewEncoder(&buf)
	// enc.Encode(msg)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// http.Handle("/", http.FileServer(http.Dir("client")))
	fmt.Println("msg", string(msg))
	w.Write(msg)
	return
}
