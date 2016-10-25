package main

import (
	"fmt"
	"time"
	"sync"
	"io/ioutil"
)

type Job struct {
	i int
	max int
	text string
}

func outputText(j *Job,goGroup *sync.WaitGroup) {

	fileName := j.text + "*.txt"
	fileContents := ""
	for j.i < j.max {
		time.Sleep(1 * time.Millisecond)
		fileContents += j.text
		fmt.Println(j.text)
		j.i++
	}
	err := ioutil.WriteFile(fileName,[]byte(fileContents),0644)
	if (err != nil) {
		panic("Something went wrong")
	}

	goGroup.Done()
}



func main(){
	goGroup := new(sync.WaitGroup)
	hello := new(Job)
	world := new(Job)

	hello.text = "hello"
	hello.i = 0
	hello.max = 3

	world.text = "World"
	world.i = 0
	world.max = 5

	go outputText(hello,goGroup)
	go outputText(world,goGroup)

	goGroup.Add(2)
	goGroup.Wait()

}


