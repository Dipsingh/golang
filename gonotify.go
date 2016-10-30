package main

import (
	"github.com/howeyc/fsnotify"
	"fmt"
	"log"
)

func main(){

	scriptDone := make(chan bool)
	dirspy,err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func (){
		for {
			select{
			case fileChange := <- dirspy.Event:
				log.Println("Something happend to file",fileChange)
			case err := <- dirspy.Error:
				log.Println("Error with Fsnotify:",err)
			}
		}
	}()
	err = dirspy.Watch("/Users/dipsingh/GOlangProjects/src/github.com/dipsingh/goconcurrency")
	if err != nil {
		fmt.Println(err)
	}
	<-scriptDone
	dirspy.Close()



}