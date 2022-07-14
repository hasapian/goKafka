package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	topic             = "events"
	broker1Address    = "localhost:9092"
	broker2Address    = "localhost:9093"
	broker3Address    = "localhost:9094"
	EVENTS_TO_PRODUCE = 5
	MESSAGE_DATA      = "this is event "
)

func TestSingleBroker(t *testing.T) {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking

	brokerAddresses := []string{broker1Address}

	go produce(ctx, brokerAddresses, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE, len(brokerAddresses))
	result := consume(ctx, brokerAddresses, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE, t)

	fmt.Println(result)

	for i := 0; i < EVENTS_TO_PRODUCE; i++ {
		j := strconv.Itoa(i)
		require.Equal(t, result[j], MESSAGE_DATA+j)
	}
}

func TestMultipleBrokers(t *testing.T) {
	// create a new context
	ctx := context.Background()

	var wg sync.WaitGroup

	brokerAddresses := []string{broker1Address, broker2Address, broker3Address}

	var chans [3]chan map[string]string
	for x := range chans {
		chans[x] = make(chan map[string]string)
	}

	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	go produce(ctx, brokerAddresses, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE, len(brokerAddresses))

	for consumerId := 0; consumerId < 3; consumerId++ {
		fmt.Println("Main: Starting consumer", consumerId)
		wg.Add(1)
		go consumeMultiple(&wg, ctx, brokerAddresses, topic, consumerId, t, chans[consumerId])
	}

	fmt.Println("Main: Waiting for consumers to finish")

	var consumerMessages [3]map[string]string
	consumerMessages[0] = <-chans[0]
	consumerMessages[1] = <-chans[1]
	consumerMessages[2] = <-chans[2]

	allMessages := make(map[string]string)

	for k := 0; k < 3; k++ {
		close(chans[k])
		for finalKey, finalValue := range consumerMessages[k] {
			allMessages[finalKey] = finalValue
		}
	}

	wg.Wait()

	fmt.Println(allMessages)
	for i := 0; i < EVENTS_TO_PRODUCE; i++ {
		j := strconv.Itoa(i)
		require.Equal(t, allMessages[j], MESSAGE_DATA+j)
	}

}
