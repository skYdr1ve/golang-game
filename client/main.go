package main

func main() {
	tcpClient := NewTcpClient("localhost:8088")
	tcpClient.Start()
}

