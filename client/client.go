package main

import (
"fmt"
"log"
"net"
"time"
"github.com/golang-game/game"
)

type TcpClient struct {
	Address string
	TcpConn *net.TCPConn
}

type Config int

var configuration Config

//Create new client
func NewTcpClient(address string) *TcpClient {
	return &TcpClient{
		Address: address,
	}
}

//Client startup function
//Call the connection function
//Call game
func (tcpClient *TcpClient) Start() {
	//Check whether connected or not
	//if not then exit
	if !tcpClient.Connect() {
		fmt.Println("Could not connect to the server")
		return
	}
	var choice int
	for {
		//Call the function that expects
		//the player to connect and
		//information for the game
		//Call the  messaging function between
		//the server and the client
		tcpClient.WaitingConfig()
		tcpClient.Receive()
		//A block of code that asks the user
		//if  he wants to continue and cheks
		//the input for correctness
		fmt.Printf("Continue playing?(1-yes,2-no): ")
		for {
			fmt.Scan(&choice)
			if choice == 1 {
				tcpClient.TcpConn.Write([]byte{1})
				break
			} else if choice == 2 {
				return
			} else {
				fmt.Println("Wrong choice. Try again: ")
			}
		}
	}
}

//Trying to connection to the server
//if after 15 seconds we did not connect to the server
//interrupt the connection attempt and exit
func (tcpClient *TcpClient) Connect() bool {
	//Check the address and port for corrctness
	//for the protocol TCP
	tcpAddress, err := net.ResolveTCPAddr("tcp", tcpClient.Address)
	if err != nil {
		log.Fatalf("Could not create TCP address from: %v %v\n", tcpClient.Address, err.Error())
	}

	connected := false
	tempTime := time.Now().Second()
	for !connected {
		log.Printf("Connecting to: %v ...\n", tcpClient.Address)
		//Trying to connect
		tcpClient.TcpConn, err = net.DialTCP("tcp4", nil, tcpAddress)
		//If not connected then try again
		if err != nil {
			log.Printf("Could not create TCP connection: %v\n", err.Error())
			time.Sleep(3 * time.Second)
			if (time.Now().Second() - tempTime) > 15 {
				return false
			}
			continue
		}

		connected = true
	}
	return true
}

//Receive information from the server
//whetere the enemy is connected and for whom client
func (tcpClient *TcpClient) WaitingConfig() {
	bytes := [1]byte{}
	for {
		//In 0 bytes is the system information
		//allowing to determine whe we are playing
		//for and whether the player is connected
		_, err := tcpClient.TcpConn.Read(bytes[0:])
		if err != nil {
			log.Printf("Could not read message: %v\n", err.Error())
			tcpClient.TcpConn.Close()
			tcpClient.Connect()
		}

		if bytes[0] == 0 {
			fmt.Println("Player search")
			continue
		} else if bytes[0] == 1 {
			configuration = 1
			fmt.Println("You are playing with crosses")
			break
		} else {
			configuration = 2
			fmt.Println("You are playing with zeroes")
			break
		}
	}
}

//Message exchange between client and server
//Get the map tic tac toe
//Check on the game situation
//Sending a message to the server
func (tcpClient *TcpClient) Receive() {
	bytes := make([]byte, game.FieldSizeInBytes+1)
	msg := make([]byte, 1)
	for {
		//Read the byte array received from the server
		//In the first 9 bytes there is a fiald for tic tac toe.
		//In 10 bytes is the system information allowing to
		//determine the course of the game
		_, err := tcpClient.TcpConn.Read(bytes[0:])

		if err != nil {
			log.Printf("Could not read message: %v\n", err.Error())
			tcpClient.TcpConn.Close()
			tcpClient.Connect()
		}
		//Call the function where we pass the first 9 bytes
		//to unsubscribe the card
		game.DrawMap(bytes[:game.FieldSizeInBytes])
		//Check the game situation which was 10 bytes
		switch game.State(bytes[game.FieldSizeInBytes]) {
		case game.GOINGON:
			break
		case game.DRAW:
			fmt.Println("Draw")
			return
		case game.PLAYER1WON:
			if configuration == 1 {
				fmt.Println("You won")
				return
			}
			fmt.Println("You lost")
			return
		case game.PLAYER2WON:
			if configuration == 2 {
				fmt.Println("You won")
				return
			}
			fmt.Println("You lost")
			return
		case game.DISCONNECTED:
			fmt.Println("Second player left")
			return
		}
		//Call the move function in which we pass
		//the field to check for correct input
		msg[0] = readPlayerInputAndCheckIt(bytes)
		_, err = tcpClient.TcpConn.Write(msg)
		if err != nil {
			log.Printf("Could not send message: %v\n", err.Error())
			return
		}

		fmt.Println("Awaiting second player's turn...")
	}
}

//The client makes a move that is checked for correctness
func readPlayerInputAndCheckIt(bytes []byte) byte {
	var x, y int
	fmt.Printf("Your turn: ")
	fmt.Scanln(&x, &y)
	for {
		if game.Check(bytes[:game.FieldSizeInBytes], x, y) {
			return byte((x-1)*3 + y - 1)
		}
		fmt.Printf("Wrong choice. Try again: ")
		fmt.Scanln(&x, &y)
	}
}