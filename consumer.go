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

	//l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
		GroupID: "my-group",
		//Logger: l,
	})

	result := make(map[string]string)
	previousOffset := -1
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		fmt.Println("received: ", string(msg.Value))
		if string(msg.Value) == "stop" {
			break
		}

		result[string(msg.Key)] = string(msg.Value)

		require.Greater(t, msg.Offset, int64(previousOffset))
		previousOffset = int(msg.Offset)
	}
	return result
}

func consumeMultiple(wg *sync.WaitGroup, ctx context.Context, brokerAddresses []string, topic string, consumerId int, t *testing.T, result chan map[string]string) {

	defer wg.Done()

	fmt.Printf("Consumer %v: Started\n", consumerId)

	//l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddresses,
		Topic:   topic,
		GroupID: "my-group",
		//Logger: l,
	})

	tempResult := make(map[string]string)

	var previousOffsets = [3]int{-1, -1, -1}

	for {
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
		fmt.Printf("received: %v\t\t", string(msg.Value))
		fmt.Printf("Consumer: %v - Partition: %v - Offset: %v\n", consumerId, msg.Partition, msg.Offset)

		require.Greater(t, msg.Offset, int64(previousOffsets[msg.Partition]))
		previousOffsets[msg.Partition] = int(msg.Offset)
	}

	fmt.Println(tempResult)

	result <- tempResult
	return
}
