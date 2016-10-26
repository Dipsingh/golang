package main

import (
	"time"
	"fmt"
)

func main(){

	ourCh := make(chan string,1)

	/*
	go func(){
	}()*/

	select {
	case <- time.After(10 * time.Second):
		fmt.Println("Enough's enough")
		close(ourCh)
	}

}
