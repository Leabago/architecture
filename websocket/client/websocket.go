package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	done := make(chan struct{})

	go func() {
		for {
			fmt.Print(">>")

			text, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}

			log.Println("send : ", text)

			err = conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				fmt.Println(err)
				return
			}

		}
	}()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("message from server: ", string(message))
		}

	}()

	for {
		select {
		case <-done:
			{
				return
			}
		case <-interrupt:
			{
				log.Println("interrupt")
				err := conn.WriteMessage(websocket.CloseNormalClosure, []byte(""))
				if err != nil {
					log.Println("write close:", err)
					return
				}

				return
			}
		}
	}

}
