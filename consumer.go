package main

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func consume(ctx context.Context, brokerAddresses []string, topic, messageValue string, numberOfMessagesToConsume int, t *testing.T) map[string]string {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	//l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
		GroupID: "my-group",
		// assign the logger to the reader
		//Logger: l,
	})

	result := make(map[string]string)
	previousOffset := -1
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		fmt.Println("received: ", string(msg.Value))
		if string(msg.Value) == "stop" {
			break
		}

		result[string(msg.Key)] = string(msg.Value)
		// after receiving the message, log its value

		require.Greater(t, msg.Offset, int64(previousOffset))
		previousOffset = int(msg.Offset)
	}
	return result
}

func consumeMultiple(wg *sync.WaitGroup, ctx context.Context, brokerAddresses []string, topic string, consumerId int, t *testing.T, result chan map[string]string) {
	//func consume(ctx context.Context, id int, t *testing.T, result chan map[string]string) {

	defer wg.Done()

	fmt.Printf("Consumer %v: Started\n", consumerId)

	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	//l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
		GroupID: "my-group",
		// assign the logger to the reader
		//Logger: l,
	})

	tempResult := make(map[string]string)

	//previousOffset := -1
	var previousOffsets = [3]int{-1, -1, -1}

	//for j := 0; j < MESSAGES_TO_PRODUCE; j++ {
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		if string(msg.Value) == "stop" {
			fmt.Println("received stop")
			fmt.Printf("Consumer %v finished\n", consumerId)
			break
		}
		tempResult[string(msg.Key)] = string(msg.Value)
		// after receiving the message, log its value
		fmt.Printf("received: %v\n", string(msg.Value))
		fmt.Printf("Consumer: %v - Partition: %v - Offset: %v\n", consumerId, msg.Partition, msg.Offset)

		require.Greater(t, msg.Offset, int64(previousOffsets[msg.Partition]))
		previousOffsets[msg.Partition] = int(msg.Offset)
	}

	fmt.Println("This is the consumer")
	fmt.Println(tempResult)

	result <- tempResult
	return
}
