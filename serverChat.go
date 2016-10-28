package main


import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var connectionCount int
var messagePool chan(string)
var Users []*User

const (
	INPUT_BUFFER_LENGTH = 140
)

type User struct {
	Name string
	ID int
	Initiated bool
	UChannel chan []byte
	Connection *net.Conn
}

func (u *User) Listen() {
	fmt.Println("Listening for",u.Name)
	for {
		select {
		case msg := <- u.UChannel:
			fmt.Println("Sending new Message to",u.Name)
			fmt.Fprintln(*u.Connection,string(msg))
		}
	}
}

type ConnectionManager struct {
	name string
	initiated bool
}

func Initiate() *ConnectionManager {
	cM := &ConnectionManager{
		name : "Chat Server 1.0",
		initiated: false,
	}
	return cM
}

func evalMessageRecipient (msg []byte, uName string) bool {
	eval := true
	expression := "@"
	re,err := regexp.MatchString(expression,string(msg))
	if err != nil {
		fmt.Println("Error Occured",err)
	}
	if re == true {
		eval = false
		pmExpression := "@" + uName
		pmRe,pmErr := regexp.MatchString(pmExpression,string(msg))
		if pmErr != nil {

		}
		if (pmRe == true){
			eval = true
		}
	}
	return eval
}

func (cM *ConnectionManager) Listen(listener net.Listener) {
	fmt.Println(cM.name, "Started")
	for {
		conn,err := listener.Accept()
		if err != nil {

		}
		connectionCount++
		fmt.Println(conn.RemoteAddr(),"connected")
		user := User{Name:"anonymous",ID:0,Initiated:false}
		Users = append(Users,&user)
		for _,u := range Users {
			fmt.Println("User Online",u.Name)
		}
		fmt.Println(connectionCount,"connection active")
		go cM.messageReady(conn,&user)
	}

}

func (cm *ConnectionManager) messageReady(conn net.Conn,user *User) {
	uChan := make(chan []byte)

	for {
		buf := make([]byte,INPUT_BUFFER_LENGTH)
		n,err := conn.Read(buf)
		if err != nil {
			conn.Close()
		}
		if n == 0 {
			conn.Close()
		}
		fmt.Println(n,"Character message from user",user.Name)
		if user.Initiated == false {
			fmt.Println("New User is",string(buf))
			user.Initiated = true
			user.UChannel = uChan
			user.Name = string(buf[:n])
			user.Connection = &conn
			go user.Listen()

			minusYouCount := strconv.FormatInt(int64(connectionCount-1),10)
			conn.Write([]byte("Welcome to the Chat,"+user.Name+". there are "+minusYouCount+ "other users"))

		} else {
			sendMessage := []byte(user.Name + ": "+ strings.TrimRight(string(buf),"\t\r\n"))

			for _,u := range Users {
				if evalMessageRecipient(sendMessage,u.Name) == true {
					u.UChannel <- sendMessage
				}
			}

		}

	}
}

func main() {
	connectionCount = 0
	serverClosed := make(chan bool)

	listener,err := net.Listen("tcp",":9000")
	if err != nil  {

	}
	connManage := Initiate()
	go connManage.Listen(listener)

	<- serverClosed
}