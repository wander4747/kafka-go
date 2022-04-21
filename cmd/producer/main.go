package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	deliveryChan := make(chan kafka.Event)

	producer := NewKafkaProducer()
	Publish("Hello Kafka consumer", "test", producer, []byte("key1"), deliveryChan)

	go DeliveryReport(deliveryChan)

	producer.Flush(1000)
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "container-kafka:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}

	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Printf(err.Error())
	}

	return p
}

func Publish(msg, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   key,
		Value: []byte(msg),
	}

	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch msg := e.(type) {
		case *kafka.Message:
			if msg.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", msg.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to topic %s [%d] at offset %v\n",
					*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
				log.Printf("Message: %s\n", msg.Value)
			}
		}
	}
}
