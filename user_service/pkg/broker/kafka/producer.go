package kafka

import (
	"context"
	"github.com/pavel/user_service/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"time"
)

type kafkaBrokerWriter struct {
	writer *kafka.Writer
}

func InitKafkaBrokerWriter(cfg config.Config) *kafkaBrokerWriter {
	kafkaConfig := kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test",
		Dialer: &kafka.Dialer{
			Timeout:  1 * time.Second,
			ClientID: "pusher",
		},
		WriteTimeout:     5 * time.Second,
		ReadTimeout:      5 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}

	return &kafkaBrokerWriter{
		writer: kafka.NewWriter(kafkaConfig),
	}
}

func (k kafkaBrokerWriter) Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	return k.writer.WriteMessages(parent, message)
}

func (k kafkaBrokerWriter) Close() error {
	return k.writer.Close()
}
