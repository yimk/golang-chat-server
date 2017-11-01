package chatroom

import (
	"fmt"
	"net"
	"strings"
	"strconv"
)

var JOIN_CHATROOM_RESPONSE_PROTOCOL = [5]string{"JOINED_CHATROOM", "SERVER_IP", "PORT", "ROOM_REF", "JOIN_ID"}
var SEND_MESSAGE_REPSONSE_PROTOCOL = [3]string{"CHAT", "CLIENT_NAME", "MESSAGE"}
var LEAVE_CHATROOM_RESPONSE_PROTOCOL = [2]string{"LEFT_CHATROOM", "JOIN_ID"}


// list of all chatRoomsRef , each chatroom's index in the list is it's reference
var chatRoomsRef = make(map[string]int)

// map of all users names, it's index is this user's join id
var userNamesJoinId = make(map[string]int)

// list of all users' corresponding group,
var usersCorrespondingGroup = make(map[string][]int)

// list of users' connection e.g users_conns[0].send(message)
var usersConns = make(map[string]net.Conn)


func RequestJoinChatroom(request string, clientConn net.Conn, port string) bool{

	// Parse request for essential information
	var requestLines = strings.Split(request, "\n")
	var roomName = strings.Split(requestLines[0], ":")[1]
	var clientName = strings.Split(requestLines[3], ":")[1]

	//Join the chatroom
	if _, exists := chatRoomsRef[roomName]; exists{
		createChatroom(roomName)
	}

	roomRef, joinId := joinChatroom(roomName, clientName, clientConn)

	//responds to the client that the joiniing is sucecesful
	returnMessage := JOIN_CHATROOM_RESPONSE_PROTOCOL[0] + ":" + roomName + "\n"
	returnMessage = returnMessage + JOIN_CHATROOM_RESPONSE_PROTOCOL[1] + ":" + getIpAddress() + "\n"
	returnMessage = returnMessage + JOIN_CHATROOM_RESPONSE_PROTOCOL[2] + ":" + port + "\n"
	returnMessage = returnMessage + JOIN_CHATROOM_RESPONSE_PROTOCOL[3] + ":" + roomRef + "\n"
	returnMessage = returnMessage + JOIN_CHATROOM_RESPONSE_PROTOCOL[4] + ":" + joinId + "\n"
	clientConn.Write([]byte(returnMessage))

	//broadcast the joining of new members in the chat room
	message := fmt.Sprintf("CHAT:%s\nCLIENT_NAME:%s\nMESSAGE:%s has joined this chatroom.\n\n", roomRef, clientName, clientName)
	broadCastWithinRoom(message, roomRef)

	//print "JOIN SUCCESFULLY\n"
	return true
}

func RequestLeavingChatroom(request string,clientConn net.Conn, serverPort string) bool{

	//Parse request for essential information
	requestLines := strings.Split(request, "\n")
	roomRef := strings.Split(requestLines[0], ":")[1]
	joinId :=  strings.Split(requestLines[1], ":")[1]
	clientName :=  strings.Split(requestLines[2], ":")[1]


	//responds to the client
	returnMessage := LEAVE_CHATROOM_RESPONSE_PROTOCOL[0] + ":" + roomRef + "\n"
	returnMessage = returnMessage + LEAVE_CHATROOM_RESPONSE_PROTOCOL[1] + ":" + joinId + "\n"

	fmt.Printf("Sent Response:\n" + returnMessage)

	clientConn.Write([]byte(returnMessage))
	message := fmt.Sprintf("CHAT:%s\nCLIENT_NAME:%s\nMESSAGE:%s has left chatroom.\n\n", roomRef, clientName, clientName)
	broadCastWithinRoom(message, roomRef)
	leaveChatroom(clientName, roomRef)
	return true
}


func RequestSendMessage(request string, clientConn net.Conn, serverPort string) bool {

	//Parse request for essential information
	requestLines := strings.Split(request, "\n")
	roomRef := strings.Split(requestLines[0], ":")[1]
	//joinId :=  strings.Split(requestLines[1], ":")[1]
	clientName :=  strings.Split(requestLines[2], ":")[1]
	message := strings.Split(requestLines[3], ":")[1]
	message = fmt.Sprintf("CHAT: %s\nCLIENT_NAME: %s\nMESSAGE: %s\n\n", roomRef, clientName, message)

	broadCastWithinRoom(message, roomRef)
	return true
}

func RequestDisconnect(request string, clientConn net.Conn, serverPort string) bool{

	// Parse request for essential information
	requestLines := strings.Split(request, "\n")
	clientName := strings.Split(requestLines[2], ":")[1]
	joinId := userNamesJoinId[clientName]

	fmt.Printf("Disconnection:")
	fmt.Printf("Client:", clientName)
	fmt.Printf("ID:", joinId)
	// fmt.Printf("Room:", usersCorrespondingGroup[joinId])


	for _, roomRef := range usersCorrespondingGroup[clientName] {
		message := fmt.Sprintf("CHAT:%s\nCLIENT_NAME:%s\nMESSAGE:%s has left chatroom.\n\n", roomRef, clientName, clientName)
		broadCastWithinRoom(message, strconv.Itoa(roomRef))
		leaveChatroom(strconv.Itoa(joinId), strconv.Itoa(roomRef))
	}

	return true
}


func Kill(){
	for _, c := range usersConns {
		c.Write([]byte("Killing the service!"))
		c.Close()
	}
}

func createChatroom(roomName string) {

	// add new chatromm
	fmt.Printf("Create CHATRoom", roomName)
	chatRoomsRef[roomName] = len(chatRoomsRef)
}

func joinChatroom(roomName string, userName string, clientConn net.Conn) (string, string){

	//get the ref of the chat room
	ref := chatRoomsRef[roomName]

	//add new user
	fmt.Println("User Conn:" + clientConn.LocalAddr().String() + "\n")
	joinId := userNamesJoinId[userName]

	if _, exists := usersConns[userName]; exists {

		// if user connection existis, then simply add new group to users chatroom record
		fmt.Printf("Client Conn Exists")
		usersCorrespondingGroup[userName] = append(usersCorrespondingGroup[userName], ref)

	} else if _, exists := userNamesJoinId[userName]; exists {

		// if user existis but user connection doesn't exists, then simply add new group and new connection to users chatroom record
		fmt.Printf("Client Name exists")
		usersCorrespondingGroup[userName] = append(usersCorrespondingGroup[userName], ref)
		usersConns[userName] = clientConn

	} else {

		joinId = len(usersConns)
		userNamesJoinId[userName] = len(userNamesJoinId)
		usersCorrespondingGroup[userName] = append(usersCorrespondingGroup[userName], ref)
		usersConns[userName] =  clientConn
		fmt.Printf("Add new user %s with join id %d", userName, joinId)

	}

	return strconv.Itoa(ref), strconv.Itoa(joinId)
 }

func broadCastWithinRoom(message string, roomRef string) {

	roomRefInt, _ :=  strconv.Atoi(roomRef)
	fmt.Printf("\nBroadCast Room: %d\n", roomRef )

	for uname, conn := range usersConns {
		fmt.Printf("JOIN ID: ", uname )
		fmt.Printf("Client: ", userNamesJoinId[uname] )

		for _, group := range usersCorrespondingGroup[uname] {
			
			if (roomRefInt == group) {
				conn.Write([]byte(message))
				break
			}
		}
	}
}

func leaveChatroom(clientName string, room string) {

	roomInt, _ := strconv.Atoi(room)
	var newUserGroup []int
	for ref, _ := range usersCorrespondingGroup[clientName] {
		if (ref != roomInt) {
			newUserGroup = append(newUserGroup, ref)
		}
	}
	usersCorrespondingGroup[clientName] = newUserGroup
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













