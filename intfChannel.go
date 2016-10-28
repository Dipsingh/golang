package main

import (
	"fmt"
	"strconv"
)

type Messenger interface {
	Relay() string
}

type Message struct {
	status string
}

func (m Message) Relay() string {
	return m.status
}

func alertMessages(v chan Messenger,i int,stats chan int) {
	m := new(Message)
	m.status = "Done with" + strconv.FormatInt(int64(i),10)
	v <- m
	stats <- 1
}


func main(){

	msg := make(chan Messenger)
	stats := make(chan int)
	for i:= 0; i < 10;i++ {
		go alertMessages(msg,i,stats)
	}
	i := 0

LOOP:
	for {
		select {
		case message,ok := <-msg:
			fmt.Println(message.Relay())
			fmt.Println("OK Status",ok)
		case <- stats:
			i++
			if (i ==9 ) {
				break LOOP
			}

		}

	}



}