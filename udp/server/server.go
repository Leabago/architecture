package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	port := ":8080"

	lp, err := net.ListenPacket("udp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer lp.Close()

	buffer := make([]byte, 1024)

	fmt.Printf("udp server listening on %s\n", port)

	for {
		_, addr, err := lp.ReadFrom(buffer)
		if err != nil {
			fmt.Println(err)
			continue
		}

		msg := string(buffer)
		fmt.Println("msg: ", msg)

		response := fmt.Sprintln(strings.ToUpper(msg))
		lp.WriteTo([]byte(response), addr)
	}

}
