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

func NewTcpClient(address string) *TcpClient {
	return &TcpClient{
		Address: address,
	}
}

func (tcpClient *TcpClient) Start() {
	if !tcpClient.Connect() {
		fmt.Println("Could not connect to the server")
		return
	}
	var choice int
	for {
		tcpClient.WaitingConfig()
		tcpClient.Receive()
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

func (tcpClient *TcpClient) Connect() bool {

	tcpAddress, err := net.ResolveTCPAddr("tcp", tcpClient.Address)

	if err != nil {
		log.Fatalf("Could not create TCP address from: %v %v\n", tcpClient.Address, err.Error())
	}

	connected := false
	tempTime := time.Now().Second()
	for !connected {

		log.Printf("Connecting to: %v ...\n", tcpClient.Address)
		tcpClient.TcpConn, err = net.DialTCP("tcp4", nil, tcpAddress)

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

func (tcpClient *TcpClient) WaitingConfig() {
	bytes := [1]byte{}
	for {
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

func (tcpClient *TcpClient) Receive() {
	bytes := make([]byte, game.FieldSizeInBytes+1)
	msg := make([]byte, 1)
	var x, y int
	for {
		n, err := tcpClient.TcpConn.Read(bytes[0:])

		if err != nil {
			log.Printf("Could not read message: %v\n", err.Error())
			tcpClient.TcpConn.Close()
			tcpClient.Connect()
		}
		game.DrawMap(bytes[:game.FieldSizeInBytes])
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
		fmt.Printf("Your turn: ")
		fmt.Scan(&x, &y)
		for {
			if game.Check(bytes[:game.FieldSizeInBytes], x, y) {
				msg[0] = byte((x-1)*3 + y - 1)
				break
			} else {
				fmt.Printf("Wrong choice. Try again: ")
				fmt.Scan(&x, &y)
			}
		}

		n, err = tcpClient.TcpConn.Write(msg)

		if err != nil {
			log.Printf("Could not send message: %v\n", err.Error())
			return
		}

		log.Printf("Sent %v bytes: %v\n", n, msg)
	}
}
