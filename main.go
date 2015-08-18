package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"time"
)

var channel chan string
var client *redis.Client

func Routine(channel chan string) {
	pubsub, err := client.Subscribe("mychannel")
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	for {
		msgi, _ := pubsub.ReceiveTimeout(2000 * time.Millisecond)
		switch msg := msgi.(type) {
		case *redis.Message:
			channel <- msg.Payload
		default:
			fmt.Println("Do nothing")
		}
	}
}

func Echo(ws *websocket.Conn) {
	var err error
	for {
		//pull data from channel and process
		token := <-channel
		msg := "Received:  " + token
		fmt.Println("Sending to client: " + msg)
		//push data back to cliend via web socket
		if token != "" {
			if err = websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("Can't send")
				break
			}
		}
	}
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//create channel that can contains string message
	channel = make(chan string)
	//spawn thread and hooked it with channel
	go Routine(channel)
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
