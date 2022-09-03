package kafka

import (
	"context"
	"github.com/pavel/message_service/config"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
	"time"
)

type kafkaBrokerWriter struct {
	writer *kafka.Writer
}

func InitKafkaBrokerWriter(cfg *config.Config) *kafkaBrokerWriter {
	writerConfig := kafka.WriterConfig{
		Brokers: []string{cfg.KafkaHost},
		Topic:   "sms",
		Dialer: &kafka.Dialer{
			Timeout:  1 * time.Second,
			ClientID: "pusher",
		},
		WriteTimeout:     5 * time.Second,
		ReadTimeout:      5 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	//test := kafka.Writer{
	//	Topic: "sms",
	//	Addr:  kafka.TCP(cfg.KafkaHost),
	//}
	return &kafkaBrokerWriter{
		//writer: &test,
		writer: kafka.NewWriter(writerConfig),
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
