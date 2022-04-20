package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	producer := NewKafkaProducer()
	Publish("message", "test", producer, nil)
	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "container-kafka:9092",
	}

	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Printf(err.Error())
	}

	return p
}

func Publish(msg, topic string, producer *kafka.Producer, key []byte) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: []byte(msg),
	}

	err := producer.Produce(message, nil)
	if err != nil {
		return err
	}

	return nil
}
