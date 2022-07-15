package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func produce(ctx context.Context, brokerAddresses []string, topic, messageValue string, numberOfMessagesToProduce int, numberOfExtraMessagesToStop int) {

	messagesProduced := 0

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
	})

	for messagesProduced < numberOfMessagesToProduce {

		err := kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(strconv.Itoa(messagesProduced)),
			Value: []byte(messageValue + strconv.Itoa(messagesProduced)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		fmt.Println("writes:", messagesProduced)
		messagesProduced++

		time.Sleep(time.Second)
	}

	for messagesProduced < numberOfMessagesToProduce+numberOfExtraMessagesToStop {
		err := kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(strconv.Itoa(messagesProduced)),
			Value: []byte("stop"),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}
		messagesProduced++
	}
}
