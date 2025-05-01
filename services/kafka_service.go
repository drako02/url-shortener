package services

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
)

func WriteKafkaEvent(topic string, key string, value string) {
	p := config.KafkaProducer

	var url models.URL
	// Fetch the user/owner of the url/shortCode and add to the event produced
	res := config.DB.Preload("User").Where("short_code=?", value).First(&url)
	if res.Error != nil {
		log.Printf("Error retrieving URL with short code %s: %v", value, res.Error)
		return
	}

	valueBytes, err := json.Marshal(ClickEvent{ShortCode: value, UserId: url.User.UID})
	if err != nil {
		log.Printf("Failed to serialize ClickEvent: %v", err)
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          valueBytes,
	}, nil)
}