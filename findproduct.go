package main


import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"math"
)

func main(){
	var number int
	reader := bufio.NewReader(os.Stdin)
	fmt.Scanln(&number)
	text,_:= reader.ReadString('\n')
	text = strings.TrimSpace(text)
	textSlice := strings.Split(text," ")

	answer,_ := strconv.Atoi(textSlice[0])

	if len(textSlice) == number {
		for _,i := range textSlice {
			integer,err := strconv.Atoi(i)
			if err != nil {
				fmt.Println("Error Occured: ",err)
			}

			answer = (answer * integer) % (int(math.Pow10(9))+7)
		}
	}

	fmt.Println(answer)
}


