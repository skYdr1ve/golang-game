package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type TcpClient struct {
	Address string
	TcpConn *net.TCPConn
}

type State byte
type Config int

const (
	GOINGON State = iota
	DRAW
	PLAYER1WON
	PLAYER2WON
	DISCONNECTED
)

var configuration Config = 0
var state State

func NewTcpClient(address string) *TcpClient {
	return &TcpClient{
		Address: address,
	}
}

func (tcpClient *TcpClient) Start() {
	if !tcpClient.Connect(){
		fmt.Println("Не удалось осуществить соединение")
		return
	}
	var choice int
	for  {
		tcpClient.WaitingConfig()
		tcpClient.Receive()
		for  {
			fmt.Println("Будите ещё играть(1-да,2-нет)")
			fmt.Scan(&choice)
			if choice==1{
				break
			}else if choice==2 {
				return
			} else {
				fmt.Println("Не верный ввод")
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
	tempTime  := time.Now().Second()
	for !connected {

		log.Printf("Connnecting to: %v ...", tcpClient.Address)
		tcpClient.TcpConn, err = net.DialTCP("tcp4", nil, tcpAddress)

		if err != nil {
			log.Printf("Could not create TCP connection: %v\n", err.Error())
			time.Sleep(3 * time.Second)
			fmt.Println(time.Now().Second()-tempTime)
			if (time.Now().Second()-tempTime)>15{
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
	for{
		n, err := tcpClient.TcpConn.Read(bytes[0:])
		if err != nil {
			log.Printf("Could not read message: %v\n", err.Error())
			tcpClient.TcpConn.Close()
			tcpClient.Connect()
		}
		log.Printf("Read %v bytes: %v", n, string(bytes[:n]))
		if(bytes[0]==0){
			fmt.Println("Player search")
		}else if(bytes[0]==1){
			configuration=1
			fmt.Println("Вы играете крестиками")
			break
		}else{
			configuration=2
			fmt.Println("Вы играете ноликами")
			break
		}
	}
}

func drawObj(obj int) string{
	if obj == 0 {
		return " "
	} else if obj == 1 {
		return "X"
	} else {
		return "0"
	}
}

func drawMap(bytes []byte) {
	for i := 0; i < 9; i++ {
		if i == 0 || i == 3 || i == 6 {
			fmt.Printf("( " + drawObj(int(bytes[i])))
		} else if i == 1 || i == 4 || i == 7 {
			fmt.Printf(" | " + drawObj(int(bytes[i])))
		} else {
			fmt.Printf(" | " + drawObj(int(bytes[i])) + " )")
		}
		if(i%3+1==0){
			fmt.Printf("\n")
		}
	}
}

func check(bytes []byte,i int) bool {
	if i > 8 || i<0 {
		return false
	}
	if bytes[i]==1 || bytes[i]==2 {
		return false;
	}else{
		return true;
	}
}

func (tcpClient *TcpClient) Receive() {
	bytes := [10]byte{}
	msg := make([]byte, 1)
	var x,y int
	for {
		n, err := tcpClient.TcpConn.Read(bytes[0:])

		if err != nil {
			log.Printf("Could not read message: %v\n", err.Error())
			tcpClient.TcpConn.Close()
			tcpClient.Connect()
		}
		log.Printf("Read %v bytes: %v", n, string(bytes[:n]))

		switch State(bytes[9]) {
			case GOINGON:
				break
			case DRAW:
				fmt.Println("НИЧЬЯ")
				return
			case PLAYER1WON:
				if configuration==1 {
					fmt.Println("ВЫ ВЫЙГРАЛИ")
					return
				}else{
					fmt.Println("ВЫ ПРОИГРАЛИ")
					return
				}
			case PLAYER2WON:
				if configuration==2 {
					fmt.Println("ВЫ ВЫЙГРАЛИ")
					return
				}else{
					fmt.Println("ВЫ ПРОИГРАЛИ")
					return
				}
			case DISCONNECTED:
				fmt.Println("ИГРОК ЛИВНУЛ")
				return
		}
		drawMap(bytes[:9])
		fmt.Println("Выбирите куда ходить")
		fmt.Scan(&x,&y)
		for {
			if check(bytes[:9],x*y-1)==true{
				msg[0]= byte(x*y - 1)
				break
			}else{
				fmt.Println("Выбирите куда ходить")
				fmt.Scan(&x,&y)
			}
		}

		_, err = tcpClient.TcpConn.Write(msg)

		if err != nil {
			log.Printf("Could not send message: %v\n", err.Error())
			return
		}

		log.Printf("Sent %v bytes\n", n)
	}
}