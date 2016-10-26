package main

import (
	"fmt"
	"sync"
	"io/ioutil"
	"strconv"
)

var writer chan bool
var rwLock sync.RWMutex


func writeFile(i int){
	fmt.Println("In Go Routine",i)
	rwLock.RLock();
	ioutil.WriteFile("test.txt",[]byte(strconv.FormatInt(int64(i),10)),0x777)
	rwLock.RUnlock();

	writer <- true

}

func main(){
	writer = make(chan bool)

	for i:=0;i<100;i++ {
		go writeFile(i)
	}

	<- writer
	fmt.Println("Done!")


}