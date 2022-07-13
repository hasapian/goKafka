package main

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	topic             = "events"
	brokerAddress     = "localhost:9092"
	EVENTS_TO_PRODUCE = 5
	MESSAGE_DATA      = "this is event "
)

func TestSingleBroker(t *testing.T) {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	go produce(ctx, brokerAddress, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE)
	result := consume(ctx, brokerAddress, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE, t)

	fmt.Println(result)

	for i := 0; i < EVENTS_TO_PRODUCE; i++ {
		j := strconv.Itoa(i)
		require.Equal(t, result[j], MESSAGE_DATA+j)
	}
}
