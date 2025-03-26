package config

import (
	"fmt"
	"os"
	// "os/signal"

	// "syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var KafkaProducer *kafka.Producer

func InitKafkaProducer() error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_URL"),

		"acks": "all",
	})
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	KafkaProducer = p

	go handleDeliveryReports(p)
	// setupGracefulShutdown(p)

	return nil
}

func handleDeliveryReports(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Produced event to topic %s: key = %-10s value - %s\n", *ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
			}

		}
	}
}

// func setupGracefulShutdown(p *kafka.Producer) {
// 	sigchan := make(chan os.Signal, 1)
// 	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

// 	go func() {
// 		<-sigchan
// 		fmt.Println("Closing Kafka producer...")
// 		// p.Flush(1000)
// 		p.Close()
// 		fmt.Println("Kafka producer closed")
// 	}()
// }
