package main

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

const (
	topic               = "events"
	brokerAddress       = "localhost:9092"
	MESSAGES_TO_PRODUCE = 5
	MESSAGE_DATA        = "this is event "
)

func TestSingleBroker(t *testing.T) {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	go produce(ctx, brokerAddress, topic, MESSAGE_DATA, MESSAGES_TO_PRODUCE)
	result := consume(ctx, t)

	fmt.Println(result)

	for i := 0; i < MESSAGES_TO_PRODUCE; i++ {
		j := strconv.Itoa(i)
		require.Equal(t, result[j], MESSAGE_DATA+j)
	}
}

func consume(ctx context.Context, t *testing.T) map[string]string {
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

		result[string(msg.Key)] = string(msg.Value)
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

		//fmt.Printf("Offset: %v\n", msg.Offset)

		if string(msg.Value) == MESSAGE_DATA+strconv.Itoa(MESSAGES_TO_PRODUCE-1) {
			break
		}
		//fmt.Printf("Previous offset: %v\n", previousOffset)
		//fmt.Printf("Current offset: %v\n", msg.Offset)
		require.Greater(t, msg.Offset, int64(previousOffset))
		previousOffset = int(msg.Offset)
	}
	return result
}
