package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/segmentio/kafka-go"
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

	ctx := context.Background()

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

	ctx := context.Background()

	var wg sync.WaitGroup

	brokerAddresses := []string{broker1Address, broker2Address, broker3Address}

	conn, err := kafka.Dial("tcp", broker1Address)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 3,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}

	var channels [3]chan map[string]string
	for x := range channels {
		channels[x] = make(chan map[string]string)
	}

	go produce(ctx, brokerAddresses, topic, MESSAGE_DATA, EVENTS_TO_PRODUCE, len(brokerAddresses))

	for consumerId := 0; consumerId < 3; consumerId++ {
		fmt.Println("Main: Starting consumer", consumerId)
		wg.Add(1)
		go consumeMultiple(&wg, ctx, brokerAddresses, topic, consumerId, t, channels[consumerId])
	}

	fmt.Println("Main: Waiting for consumers to finish")

	var consumerMessages [3]map[string]string
	consumerMessages[0] = <-channels[0]
	consumerMessages[1] = <-channels[1]
	consumerMessages[2] = <-channels[2]

	allMessages := make(map[string]string)

	for k := 0; k < 3; k++ {
		close(channels[k])
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
