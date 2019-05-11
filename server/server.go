package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	game "github.com/golang-game/game"
)

const (
	connPort    = "5050"
	connType    = "tcp"
	connAddress = "localhost"
)

var clientCounter int32
var gameXO game.GameState

func main() {
	var servSync sync.WaitGroup
	servSync.Add(1)
	address, err := net.ResolveTCPAddr(connType, connAddress+":"+connPort)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	l, err := net.ListenTCP(connType, address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Successfully started listening on", connAddress+":"+connPort)
	var mutex sync.Mutex
	gameConcluded := make(chan int)
	for true {
		for clientCounter < 2 {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			clientCounter++
			fmt.Println("Client", clientCounter, "connected")
			go handleClient(conn, clientCounter, &mutex, &servSync, gameConcluded)
		}
		servSync.Wait()
	}
}

func handleClient(conn net.Conn, id int32, mutex *sync.Mutex, servSync *sync.WaitGroup, gameConcluded chan int) {
	defer servSync.Done()
	for true {
		for clientCounter < 2 {
			_, err := conn.Write([]byte("0"))
			if err != nil {
				atomic.AddInt32(&clientCounter, -1)
				fmt.Println("Client", id, "disconnected")
				gameXO.ResetGame()
				gameConcluded <- -1
				return
			}
			time.Sleep(time.Second)
		}
		_, err := conn.Write([]byte{byte(id)})
		if err != nil {
			atomic.AddInt32(&clientCounter, -1)
			fmt.Println("Client", id, "disconnected")
			gameXO.ResetGame()
			gameConcluded <- -1
			return
		}
		continuePlaying := true
		for continuePlaying {
			mutex.Lock()
			select {
			case exitCode := <-gameConcluded:
				if exitCode == -1 {
					gameXO.State = game.DISCONNECTED
				}
				_, err = conn.Write(append(gameXO.PlayingField, byte(gameXO.State)))
				if err != nil {
					atomic.AddInt32(&clientCounter, -1)
					fmt.Println("Client", id, "disconnected")
					return
				}
			default:
				turn := make([]byte, 3)
				_, err = conn.Write(append(gameXO.PlayingField, byte(gameXO.State)))
				if err != nil {
					atomic.AddInt32(&clientCounter, -1)
					fmt.Println("Client", id, "disconnected")
					gameXO.ResetGame()
					gameConcluded <- -1
					return
				}
				_, err := conn.Read(turn)
				if err != nil {
					atomic.AddInt32(&clientCounter, -1)
					fmt.Println("Client", id, "disconnected")
					gameXO.ResetGame()
					gameConcluded <- -1
					return
				}
				gameXO.PlayingField[int(turn[0])] = byte(id)
				gameXO.CheckState()
				if gameXO.State != game.GOINGON {
					continuePlaying = false
					gameConcluded <- 1
				}
			}
			mutex.Unlock()
		}
		_, err = conn.Write(append(gameXO.PlayingField, byte(gameXO.State)))
		if err != nil {
			atomic.AddInt32(&clientCounter, -1)
			fmt.Println("Client", id, "disconnected")
			return
		}
		gameXO.ResetGame()
	}
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:4040")
	if err != nil {
		log.Fatal(err)
		fmt.Println("Error while retrieving local ip address")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
