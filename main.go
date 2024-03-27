package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	handleService "golang-hotstar/services"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// func processEmoji(resW http.ResponseWriter, req *http.Request) {

// }

func main() {
	server := mux.NewRouter()
	server.HandleFunc("/emoji", handleService.HandleEmoji).Methods(http.MethodPost)
	fmt.Printf("Starting server at port 8081\n")

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "127.0.0.1:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		fmt.Println("error while creating consumer", err)
	}

	c.SubscribeTopics([]string{"my-topic"}, nil)

	go func() {
		run := true

		for run {
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}

		c.Close()
	}()

	log.Fatal(http.ListenAndServe(":8081", server))

}
