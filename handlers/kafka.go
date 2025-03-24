package handlers

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/drako02/url-shortener/config"
)

func WriteKafkaEvent (topic string, key string, value string) {
	p := config.KafkaProducer
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key: []byte(key),
		Value: []byte(value),
	}, nil)
}


