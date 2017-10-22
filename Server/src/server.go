package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"src/lib"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
	MAX_DATA_RECV = 9999
	BACKLOG = 50
)

//var chatroom = Chatroom()



func main() {

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
	
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {

  	// Make a buffer to hold incoming data.
  	buf := make([]byte, MAX_DATA_RECV)

  	// Read the incoming connection into the request.
  	reqLen, err := conn.Read(buf)
	request := string(buf[:reqLen])

  	if err != nil {
		fmt.Println("Error reading:", err.Error())
  	}

  	if(strings.Compare(request, "KILL_SERVICE\n") == 1) {

	  	fmt.Printf("Kill the service\n" )
	  	//chatroom.kill()
	  	conn.Close()

   	} else if(strings.Compare(request, "HELO text\n") == 1) {

   		fmt.Printf("Send back Hello\n" )
   		//"HELO text\nIP:[ip address]\nPort:[port number]\nStudentID:[your student ID]\n"
		ip := getIpAddress()
		returnMessage := "HELO text\nIP::" + ip + "\nPort:" + CONN_PORT + "\nStudentID:" + "13329643" + "\n"
		fmt.Printf(returnMessage)
		conn.Write([]byte(returnMessage))
		fmt.Printf(request)

   	} else if(strings.Contains(request, "JOIN_CHATROOM")) {

   		fmt.Printf("It is a JOIN CHATROOM REQUEST\n") 
		//chatroom.requestJoinChatroom(request, conn, CONN_PORT)

   	} else if(strings.Contains(request, "LEAVE_CHATROOM")) {

   		fmt.Printf("It is a LEAVE CHATROOM REQUEST\n") 
		//chatroom.requestLeavingChatroom(request, conn, CONN_PORT)

   	} else if(strings.Contains(request, "CHAT")) {

   		fmt.Printf("It is a JOIN CHATROOM REQUEST\n") 
		//chatroom.requestSendMessage(request, conn, CONN_PORT)

   	} else if(strings.Compare(request, "DISCONNECT") == 1) {

   		fmt.Printf("It is a LEAVE CHATROOM REQUEST\n") 
		//chatroom.requestDisconnect(request, conn, CONN_PORT)

   	} else {

   		fmt.Printf("Nothing interesting\n")
   		conn.Write([]byte("Nothing interesting."))

   	}

   	fmt.Printf("Task Complete")
   	conn.Close()

}

func getIpAddress() string{

	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			return ip.String()
		}
	}
	return ""
}

























