package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	game "github.com/golang-game/game"
)

const (
	connPort    = "8088"
	connType    = "tcp"
	connAddress = "localhost"
)

var clientCounter int32
var gameXO game.GameState

func main() {
	var servSync sync.WaitGroup
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
	gameConcluded := make(chan int, 1)
	gameXO = game.New()
	for true {
		servSync.Add(1)
		for clientCounter < 2 {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			atomic.AddInt32(&clientCounter, 1)
			newID := clientCounter
			fmt.Println("Client", newID, "connected")
			go handleClient(conn, newID, &mutex, &servSync, gameConcluded)
		}
		servSync.Wait()
	}
}

func handleClient(conn net.Conn, id int32, mutex *sync.Mutex, servSync *sync.WaitGroup, gameConcluded chan int) {
	defer servSync.Done()
	defer atomic.AddInt32(&clientCounter, -1)
	defer fmt.Println("Client", id, "disconnected")
	defer gameXO.ResetGame()
	for true {
		for clientCounter < 2 {
			_, err := conn.Write([]byte{0})
			if err != nil {
				return
			}
			time.Sleep(time.Second)
		}
		_, err := conn.Write([]byte{byte(id)})
		if err != nil {
			if clientCounter == 2 {
				gameConcluded <- -1
			}
			return
		}
		if id == 2 {
			waitTime := rand.Intn(700)
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
		}
		continuePlaying := true
		for continuePlaying {
			time.Sleep(time.Duration(100) * time.Millisecond)
			//log.Printf("Before lock: %d\n", id)
			mutex.Lock()
			//log.Printf("After lock: %d\n", id)
			select {
			case exitCode := <-gameConcluded:
				if exitCode == -1 {
					gameXO.State = game.DISCONNECTED
				}
				continuePlaying = false
			default:
				turn := make([]byte, 1)
				_, err = conn.Write(append(gameXO.PlayingField, byte(gameXO.State)))
				if err != nil {
					if clientCounter == 2 {
						gameConcluded <- -1
					}
					continuePlaying = false
				}
				_, err = conn.Read(turn)
				//log.Printf("Read %v bytes: %v", n, turn)
				if err != nil {
					if clientCounter == 2 {
						gameConcluded <- -1
					}
					continuePlaying = false
				}
				gameXO.PlayingField[int(turn[0])] = byte(id)
				gameXO.CheckState()
				if gameXO.State != game.GOINGON {
					continuePlaying = false
					if clientCounter == 2 {
						gameConcluded <- 1
					}
				}
			}
			mutex.Unlock()
			//log.Printf("After unlock: %d\n", id)
		}
		_, err = conn.Write(append(gameXO.PlayingField, byte(gameXO.State)))
		if err != nil {
			break
		}
		continuation := []byte{0}
		_, err = conn.Read(continuation)
		if err != nil {
			break
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
