package main

import "net"
import "fmt"
import "bufio"

func main() {

	// connect to this socket
	conn, _ := net.Dial("tcp", "127.0.0.1:3333")
	for {
		// read in input from stdin
		text := "JOIN_CHATROOM: room\nCLIENT_IP: 0\nPORT: 0\nCLIENT_NAME: Peter"

		// send to socket
		fmt.Fprintf(conn, text + "\n")

		// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}