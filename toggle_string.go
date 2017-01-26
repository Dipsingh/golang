package main

import (
	"fmt"
	"strings"
	"bytes"
)


func main(){

	var input string
	fmt.Scanln(&input)
	var buffer bytes.Buffer
	new := strings.Split(input,"")
	for _, i := range new {
		if (i == strings.ToLower(i)) {
			tmpUpper := strings.ToUpper(i)
			buffer.WriteString(tmpUpper)
		}
		if (i == strings.ToUpper(i)) {
			tmpLower := strings.ToLower(i)
			buffer.WriteString(tmpLower)

		}
	}

	fmt.Println(buffer.String())

}
