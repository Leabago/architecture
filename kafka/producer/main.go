package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	// reader := bufio.NewReader(os.Stdin)
	i := 0

	for {
		fmt.Print(">>")

		text := "send: " + strconv.Itoa(i)
		i++

		// text, err := reader.ReadString('\n')
		// if err != nil {
		// 	log.Fatal("failed to read string:", err)
		// }

		// text = strings.Trim(text, "\n")

		// if text == "EXIT" {
		// 	break
		// }

		fmt.Printf("text: {%s}\n", text)

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(text)},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		time.Sleep(1 * time.Second)

	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
