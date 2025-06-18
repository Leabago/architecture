package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">>")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(conn, text)

		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("->", message)

		if strings.TrimSpace(text) == "STOP" {
			break
		}
	}

}
