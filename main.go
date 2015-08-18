package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"code.google.com/p/go.net/websocket"
)

var channel chan string

func Routine(channel chan string) {
	for {
		time.Sleep(2000 * time.Millisecond)
		t := time.Now()
		//send data back to main process via channel
		channel <- t.Format("20060102150405")
	}
}

func Echo(ws *websocket.Conn) {
    var err error
    for {
		//pull data from channel and process
		token := <- channel
        msg := "Received:  " + token
        fmt.Println("Sending to client: " + msg)
		//push data back to cliend via web socket
        if err = websocket.Message.Send(ws, msg); err != nil {
            fmt.Println("Can't send")
            break
        }
    }
}

func main() {
	//create channel that can contains string message
	channel = make(chan string)
	//spawn thread and hooked it with channel
	go Routine(channel)
    http.Handle("/", websocket.Handler(Echo))

    if err := http.ListenAndServe(":1234", nil); err != nil {
        log.Fatal("ListenAndServe:", err)
    }
}
