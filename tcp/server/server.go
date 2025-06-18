package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	port := ":8080"

	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer l.Close()

	fmt.Printf("TCP server listening on %s\n", port)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		msg := strings.TrimSpace(string(netData))
		if msg == "STOP" {
			break
		}

		response := fmt.Sprintf("Server: %s\n", strings.ToUpper(msg))
		conn.Write([]byte(response))
	}
}
