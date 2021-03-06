# goKafka
Kafka testing in Go for single and multiple brokers. The main concept of the implementation was based on the [post](https://www.sohamkamani.com/golang/working-with-kafka/) of Soham Kamani, which was enriched with further functionalities. 

- Go client library for Kafka: [kafka-go](https://github.com/segmentio/kafka-go)
- Docker compose for single broker was created based on the [tutorial](https://developer.confluent.io/get-started/go/#kafka-setup) by Confluent
- Docker compose for multiple brokers: [kafka-stack-docker-compose](https://github.com/conduktor/kafka-stack-docker-compose)
- The [require](https://pkg.go.dev/github.com/stretchr/testify) package was used for the assertions of the tests

Each kafka message has a key and value. The key is used to decide which partition (and consequently, which broker) the message gets published on. When all the messages are produced, an extra message (or 3 if there are 3 partitions) with value **stop** is produced to stop the consumer(s).
| Key | Value            |
| --- | ---------------- |
| 0   | this is event 0  |
| 1   | this is event 1  |
| 2   | this is event 2  |
| 3   | stop             |

# Test Kafka with Single Broker
*TestSingleBroker* is a test case for Kafka with **one broker and one partition**. The gihub action yaml file for this test is *github-action-single-broker*, in which the containers and the topic are created. The [confluentinc](https://hub.docker.com/u/confluentinc) images are used for the [zookeeper](https://hub.docker.com/r/confluentinc/cp-zookeeper) and the [kafka](https://hub.docker.com/r/confluentinc/cp-kafka) broker. The topic is created with the command

`docker compose exec broker kafka-topics --bootstrap-server localhost:9092 --topic events --create --partitions 1 --replication-factor 1`

A specified number of events is created, which are consumed by a single consumer. The consumer stops when a message with value **stop** is received.

## GitHub Actions Kafka Single Broker
The workflow for the single broker has the following steps:
- Check out repository
- Pull images for broker and zookeeper via docker compose and start containers
- Create the necessary topic to use during tests
- Run test for single broker

## Steps to run locally
1. Clone repository: `git clone https://github.com/hasapian/goKafka.git`
2. Create Zookeeper and Kafka containers with docker compose (install if needed following [this](https://docs.docker.com/compose/install/) manual): `docker compose up -d`
3. Create a topic with name **events**: `docker compose exec broker kafka-topics --bootstrap-server localhost:9092 --topic events --create --partitions 1 --replication-factor 1`
4. Execute the test case: `go test -v -run TestSingle`

# Test Kafka with Multiple Brokers
*TestMultipleBrokers* is a TC for Kafka with **3 brokers, 3 partitions and replication factor 3** (transaction.state.log.min.isr has the default value: 2). In this test case, the topic (*events*) is created during the execution, using the [kafka-go](https://github.com/segmentio/kafka-go) library. 3 goroutines are created for the 3 consumers, each of which stops when a message with value **stop** is received.
- Github action yaml: *github-action-multiple-brokers*. 
- For the Zookeeper-Kafka setup the file [zk-single-kafka-multiple.yml](https://github.com/conduktor/kafka-stack-docker-compose) is used. 

## GitHub Actions Kafka Multiple Brokers

The workflow for the multiple brokers has the following steps:
- Check out repository
- Pull images for brokers and zookeeper via docker compose and start containers
- Add some delay for brokers and zookeeper to be ready 
- Run test for multiple brokers

## Steps to run locally
1. Clone repository: `git clone https://github.com/hasapian/goKafka.git`
2. Create Zookeeper and Kafka containers with docker compose : `docker-compose -f zk-single-kafka-multiple.yml up -d`
3. Execute the test case: `go test -v -run TestMultiple`
