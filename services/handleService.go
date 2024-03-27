package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type EmojiRequest struct {
	Emoji string `json:"emoji"`
}

type ProcessEmoji struct {
	pushChannel chan string
}

func (er *ProcessEmoji) pushEmoji(emoji string) {
	fmt.Println("process emoji: ", emoji)

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "127.0.0.1:9094"})

	if err != nil {
		fmt.Println("error while creating producer", err)
	}

	defer producer.Close()

	topic := "my-topic"

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(emoji),
	}, nil)

	producer.Flush(15 * 1000)
}

func HandleEmoji(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req EmojiRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "error while decoding request body", http.StatusBadRequest)
		return
	}

	processemoji := &ProcessEmoji{
		pushChannel: make(chan string),
	}

	go func() {
		processemoji.pushChannel <- req.Emoji
	}()

	go func() {
		// for {
		// 	// select {

		// 	// case response := <-processemoji.pushChannel:
		// 	// 	fmt.Println("LOGED LGOED", response)
		// 	// case processemoji.pushChannel <- req.Emoji:
		// 	// 	// fmt.Println("LOGED ÷LGOED")
		// 	// 	processemoji.pushChannel <- req.Emoji
		// 	// }

		// 	select {
		// 	case response := <-processemoji.pushChannel:
		// 		processemoji.pushEmoji(response)
		// 		go w.Write([]byte("Received emoji: " + response))
		// 	}

		// }

		for {
			select {
			case response := <-processemoji.pushChannel:
				processemoji.pushEmoji(response)
				// w.Write([]byte(response))
			}
		}
	}()

	// go func() {
	// 	response := <-processemoji.pushChannel
	// 	processemoji.pushEmoji(response)
	// }()

	response := "Received emoji: " + req.Emoji
	w.Write([]byte(response))
}
