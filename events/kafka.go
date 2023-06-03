package events

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"os"
)

var _WRITER *kafka.Writer

func Init() {

	if _WRITER == nil {
		kafkaURL := os.Getenv("KAFKA_URL")
		topic := os.Getenv("KAFKA_TOPIC")

		_WRITER = getKafkaWriter(kafkaURL, topic)
	}
}

func Close() {
	err := _WRITER.Close()
	if err != nil {
		log.Error("Error closing kafka reader", err)
	}
}

func ProduceKafka(key string, body interface{}) {
	bodyStr, _ := json.Marshal(body)
	msg := kafka.Message{
		Key:   []byte(key),
		Value: bodyStr,
	}
	err := _WRITER.WriteMessages(context.Background(), msg)

	if err != nil {
		log.Error("Unable to write to kafka", err)
	}
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
