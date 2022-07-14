package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func produce(ctx context.Context, brokerAddresses []string, topic, messageValue string, numberOfMessagesToProduce int, numberOfExtraMessagesToStop int) {
	// initialize a counter
	messagesProduced := 0

	//l := log.New(os.Stdout, "kafka writer: ", 0)
	// intialize the writer with the broker addresses, and the topic
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
		// assign the logger to the writer
		//Logger: l,
	})

	for messagesProduced < numberOfMessagesToProduce {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(messagesProduced)),
			// create an arbitrary message payload for the value
			Value: []byte(messageValue + strconv.Itoa(messagesProduced)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", messagesProduced)
		messagesProduced++
		// sleep for a second
		time.Sleep(time.Second)
	}

	for messagesProduced < numberOfMessagesToProduce+numberOfExtraMessagesToStop {
		err := kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(messagesProduced)),
			// create an arbitrary message payload for the value
			Value: []byte("stop"),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		messagesProduced++
	}
}
