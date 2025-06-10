package services

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/drako02/url-shortener/config"
	"github.com/drako02/url-shortener/models"
	"github.com/drako02/url-shortener/utils"
)

func WriteKafkaEvent(topic string, key string, value string) error {
	p := config.KafkaProducer

	var event *utils.ClickEvent
	if err := json.Unmarshal([]byte(value), &event); err != nil {
		log.Printf("Error unmarshaling value: %v", err)
		return err
	}
	var url models.URL
	// Fetch the user/owner of the url/shortCode and add to the event produced
	if err := config.DB.Preload("User").Where("short_code=?", event.ShortCode).First(&url).Error; err != nil {
		log.Printf("Error retrieving URL with short code %s: %v", event.ShortCode, err)
		return err
	}
	// Set the UserID in the ClickEvent from the URL's User data
	event.UserId = url.UserId
	valueBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to serialize ClickEvent: %v", err)
		return err
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          valueBytes,
	}, nil)

	return nil
}
