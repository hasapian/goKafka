package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func consume(ctx context.Context, brokerAddress, topic, messageValue string, numberOfMessagesToConsume int, t *testing.T) map[string]string {
	// create a new logger that outputs to stdout
	// and has the `kafka reader` prefix
	//l := log.New(os.Stdout, "kafka reader: ", 0)
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
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
