package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	//play safe. wait for kafka container start (only for docker-ized env)
	time.Sleep(10 * time.Second)

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		log.Fatal("missing KAFKA_BROKER env variable")
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		log.Fatal("missing KAFKA_TOPIC env variable")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatal("unable to create kafka producer", err)
	}

	defer func() {
		p.Close()
		log.Println("kafka producer closed")
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	close := false

	go func() {
		log.Println("starting producer loop...")
		for !close {
			stationID := "station-" + strconv.Itoa(rand.Intn(9)+1)
			temp := strconv.Itoa(rand.Intn(10) + 60)
			state := "state-" + strconv.Itoa(rand.Intn(19)+1)
			created := time.Now().Format(time.RFC3339)

			weather := Weather{StationID: stationID, Temperature: temp, State: state, Created: created}

			weatherMsg, err := json.Marshal(weather)
			if err != nil {
				log.Println("unable to marshal weather to JSON payload", err)
			}

			log.Println("weather payload", string(weatherMsg))

			msg := &kafka.Message{Key: []byte(stationID), Value: weatherMsg, TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}}

			//we don't care about delivery status/reports (hence using nil chan and neither are we tracking the p.Events() channel)
			err = p.Produce(msg, nil)
			if err != nil {
				log.Println("failed to enqueue message", err)
			}

			//take it easy!
			time.Sleep(3 * time.Second)
		}
		log.Println("exited producer loop...")
	}()

	log.Println("press ctrl+c to exit")
	<-exit
	log.Println("exit signalled")

	close = true
	//flush all messages in the internal queue before shutting down
	p.Flush(5000)
}

//Weather data
type Weather struct {
	StationID   string `json:"stationid"`
	Temperature string `json:"temp"`
	State       string `json:"state"`
	Created     string `json:"created"`
}
